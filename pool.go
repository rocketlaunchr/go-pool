package pool

import (
	"context"
	// "fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// Options configures the Pool struct.
type Options struct {
	// Initial creates an initial number of ready-to-use items in the pool.
	Initial int
	// Max sets the maximum number of items kept in the pool.
	Max *int
	// DisableCount, when set, disables the pool's Count function.
	// Only set this if you need a runtime Finalizer for the item returned
	// by the factory (alternatively, wrap your item in another struct, with the Finalizer
	// added to the original item).
	DisableCount bool
}

// New creates a new Pool.
// opts accepts either an int (representing the max) or an Options struct.
func New(opts ...interface{}) Pool {
	if len(opts) == 0 {
		return Pool{}
	}
	pool := Pool{}
	switch o := opts[0].(type) {
	case int:
		pool.semMax = semaphore.NewWeighted(int64(o))
		return pool
	case Options:
		// max
		if o.Max != nil {
			if o.Initial > *o.Max {
				panic("Initial must not exceed Max")
			}
			pool.semMax = semaphore.NewWeighted(int64(*o.Max))
		}

		// initial
		if o.Initial > 0 {
			pool.initial = &o.Initial
		}

		// noCount
		pool.noCount = o.DisableCount

		return pool
	}

	panic("opts must be an int or Options struct")
}

// A Pool is a set of temporary objects that may be individually saved and
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
// On the other hand, a free list maintained as part of a short-lived object is
// not a suitable use for a Pool, since the overhead does not amortize well in
// that scenario. It is more efficient to have such objects implement their own
// free list.
//
// A Pool must not be copied after first use.
type Pool struct {
	noCopy noCopy

	initial *int // if nil, then initial items have already been created (or initial option was no set)

	syncPool sync.Pool
	semMax   *semaphore.Weighted

	noCount bool
	count   uint32 // count keeps track of approximately how many items are in the pool
}

// SetFactory specifies a function to generate an item when Borrow is called.
// It must not be called concurrently with calls to Borrow.
//
// NOTE: factory should generally only return pointer types, since a pointer can be put into the return interface
// value without an allocation.
func (p *Pool) SetFactory(factory func() interface{}) {

	p.syncPool.New = func() interface{} {
		newItem := factory()

		if !p.noCount {
			atomic.AddUint32(&p.count, 1) // p.count++
			runtime.SetFinalizer(newItem, func(newItem interface{}) {
				atomic.AddUint32(&p.count, ^uint32(0)) // p.count--
				// fmt.Printf("Factory Item has been garbage collected. (%d left)\n", p.count)
			})
		}

		// fmt.Printf("New Factory Item created (%d in pool)\n", p.count)
		return newItem
	}

	if p.initial != nil {
		// create initial number of items
		items := []interface{}{}

		// create new items
		for i := 0; i < *p.initial; i++ {
			items = append(items, p.borrow())
		}
		// return new items
		for j := len(items) - 1; j >= 0; j-- {
			p.returnItem(items[j])
		}
		p.initial = nil
	}
}

func (p *Pool) borrow() interface{} {
	if p.semMax != nil {
		p.semMax.Acquire(context.Background(), 1)
	}
	wrap := itemWrapPool.Get().(*ItemWrap)
	item := p.syncPool.Get()

	wrap.Item = item
	wrap.pool = p
	return wrap
}

func (p *Pool) returnItem(x interface{}) {
	wrap := x.(*ItemWrap)
	wrap.pool = nil
	if wrap.invalid {
		wrap.invalid = false
	} else {
		p.syncPool.Put(wrap.Item)
	}
	wrap.Item = nil
	itemWrapPool.Put(wrap)
	if p.semMax != nil {
		p.semMax.Release(1)
	}
}

// Borrow obtains an item from the pool.
// If the Max option is set, then this function
// will block until an item is returned back into the pool.
//
// After the item is no longer required, you must call
// Return on the item.
func (p *Pool) Borrow() *ItemWrap {
	return p.borrow().(*ItemWrap)
}

// ReturnItem returns an item back to the pool.
// Usually this function is never called, as the recommended
// approach is to call Return on the item.
func (p *Pool) ReturnItem(x *ItemWrap) {
	p.returnItem(x)
}

// Count returns approximately the number of items in the pool (idle and in-use).
// If you want an accurate number, call runtime.GC() twice before calling Count (not recommended).
func (p *Pool) Count() int {
	if p.noCount {
		return 0
	}
	return int(atomic.LoadUint32(&p.count))
}
