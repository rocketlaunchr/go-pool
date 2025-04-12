package pool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// Options configures the Pool struct.
type Options struct {
	// Initial creates an initial number of ready-to-use items in the pool.
	Initial uint32

	// Max represents the maximum number of items you can borrow. This prevents
	// unbounded growth in the pool.
	//
	// Depending on the timing of Returns and Factory calls, the maximum number of
	// items in the pool can exceed Max by a small number for a short time.
	Max uint32

	// EnableCount, when set, enables the pool's Count function.
	//
	// NOTE: If you set this AND you need to set your own runtime Finalizer on the item,
	// wrap your item in another struct, with the Finalizer added to the inner item.
	EnableCount bool
}

// A Pool is a set of temporary objects that may be individually stored and
// retrieved.
//
// Any item stored in the Pool may be removed automatically at any time without
// notification. If the Pool holds the only reference when this happens, the
// item might be deallocated.
//
// A Pool is safe for use by multiple goroutines simultaneously.
//
// Pool's purpose is to cache allocated but unused items for later reuse,
// relieving pressure on the garbage collector. That is, it makes it easy to
// build efficient, thread-safe free lists. However, it is not suitable for all
// free lists.
//
// An appropriate use of a Pool is to manage a group of temporary items
// silently shared among and potentially reused by concurrent independent
// clients of a package. Pool provides a way to amortize allocation overhead
// across many clients.
//
// The Pool scales under load and shrinks when quiescent.
//
// On the other hand, a free list maintained as part of a short-lived object is
// not a suitable use for a Pool, since the overhead does not amortize well in
// that scenario. It is more efficient to have such objects implement their own
// free list.
//
// A Pool must not be copied after first use.
type Pool[T any] struct {
	noCopy noCopy

	initial *uint32 // if nil, then initial items have already been created (or initial option was no set)

	itemWrapPool *sync.Pool
	syncPool     sync.Pool
	semMax       *semaphore.Weighted

	enableCount      bool
	count            uint32 // count keeps track of approximately how many items are in existence (in the pool and in-use).
	countBorrowedOut uint32
}

// New creates a new Pool.
// opts accepts either an int (representing the max) or an Options struct.
func New[T any](opts ...any) Pool[T] {
	pool := Pool[T]{
		itemWrapPool: &sync.Pool{
			New: func() any { return new(ItemWrap[T]) },
		},
	}
	if len(opts) == 0 {
		return pool
	}
	switch o := opts[0].(type) {
	case int:
		pool.semMax = semaphore.NewWeighted(int64(o))
	case uint32:
		pool.semMax = semaphore.NewWeighted(int64(o))
	case Options:
		// max
		if o.Max != 0 {
			if o.Initial > o.Max {
				panic("Initial must not exceed Max")
			}
			pool.semMax = semaphore.NewWeighted(int64(o.Max))
		}

		// initial
		if o.Initial > 0 {
			pool.initial = &o.Initial
		}

		// enableCount
		pool.enableCount = o.EnableCount
	default:
		panic("opts must be an int or Options struct")
	}
	return pool
}

// SetFactory specifies a function to generate a new item when Borrow is called.
// It must not be called concurrently with calls to Borrow.
//
// NOTE: factory should generally only return pointer types, since a pointer can be put into the return interface
// value without an allocation.
func (p *Pool[T]) SetFactory(factory func() T) {

	p.syncPool.New = func() any {
		newItem := factory()

		if p.enableCount {
			atomic.AddUint32(&p.count, 1) // p.count++
			runtime.SetFinalizer(newItem, func(newItem any) {
				atomic.AddUint32(&p.count, ^uint32(0)) // p.count--
				// fmt.Printf("+++Factory Item has been garbage collected. (%d left)\n", p.count)
			})
		}

		// fmt.Printf("***New Factory Item created (%d in pool)\n", p.count)
		return newItem
	}

	if p.initial != nil {
		// create initial number of items
		items := make([]*ItemWrap[T], 0, *p.initial)

		// create new items
		for i := uint32(0); i < *p.initial; i++ {
			items = append(items, p.borrow())
		}
		// return new items
		for j := len(items) - 1; j >= 0; j-- {
			p.returnItem(items[j])
		}
		p.initial = nil
	}
}

func (p *Pool[T]) borrow() *ItemWrap[T] {
	if p.semMax != nil {
		p.semMax.Acquire(context.Background(), 1)
	}
	if p.enableCount {
		atomic.AddUint32(&p.countBorrowedOut, 1) // p.countBorrowedOut++
	}

	wrap := p.itemWrapPool.Get().(*ItemWrap[T])
	item := p.syncPool.Get()

	wrap.Item = item.(T)
	wrap.pool = p
	return wrap
}

func (p *Pool[T]) returnItem(wrap *ItemWrap[T]) {
	if !wrap.invalid {
		p.syncPool.Put(wrap.Item)
	}
	wrap.Reset()

	p.itemWrapPool.Put(wrap)
	if p.enableCount {
		atomic.AddUint32(&p.countBorrowedOut, ^uint32(0)) // p.countBorrowedOut--
	}
	if p.semMax != nil {
		p.semMax.Release(1)
	}
}

// Borrow obtains an item from the pool.
// If the Max option is set, then this function will
// block until an item is returned back into the pool.
//
// After the item is no longer required, you must call
// Return on the item.
func (p *Pool[T]) Borrow() *ItemWrap[T] {
	return p.borrow()
}

// ReturnItem returns an item back to the pool. Usually
// developer's never call this function, as the recommended
// approach is to call Return on the item.
func (p *Pool[T]) ReturnItem(x *ItemWrap[T]) {
	p.returnItem(x)
}

// Count returns approximately the number of items in the pool (idle).
// If you want an accurate number, call runtime.GC() twice before calling Count (not recommended).
//
// NOTE: Count can exceed both the Initial and Max value by a small number for a short time.
func (p *Pool[T]) Count() uint32 {
	if !p.enableCount {
		return 0
	}
	c := atomic.LoadUint32(&p.count)
	b := atomic.LoadUint32(&p.countBorrowedOut)
	if c > b {
		return c - b
	}
	return 0
}

// OnLoan returns how many items are in-use.
func (p *Pool[T]) OnLoan() uint32 {
	if !p.enableCount {
		return 0
	}
	return atomic.LoadUint32(&p.countBorrowedOut)
}
