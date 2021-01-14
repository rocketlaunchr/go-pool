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
	// Max sets the maximum number items kept in the pool.
	Max *int
}

// New creates a new Pool.
// opts accepts either an int (representing the max) or an Options struct.
func New(opts ...interface{}) Pool {
	if len(opts) == 0 {
		return Pool{}
	}

	switch o := opts[0].(type) {
	case int:
		pool := Pool{}
		pool.semMax = semaphore.NewWeighted(int64(o))
		return pool
	case Options:
		pool := Pool{}

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

	count uint32 // count keeps track of approximately how many items are in the pool
}

// SetFactory specifies a function to generate an item when GetItem is called.
// It must not be called concurrently with calls to GetItem.
func (p *Pool) SetFactory(factory func() interface{}) {

	p.syncPool.New = func() interface{} {
		defer func() { atomic.AddUint32(&p.count, 1) }()
		newItem := factory()

		runtime.SetFinalizer(newItem, func(newItem interface{}) {
			atomic.AddUint32(&p.count, ^uint32(0)) // p.count--
			// fmt.Printf("Factory Item [%d] has been garbage collected. (%d left)\n", i, count)
		})

		if p.semMax != nil {
			p.semMax.Acquire(context.Background(), 1)
		}
		// fmt.Printf("New Factory Item [%d] created (%d in pool)\n", i, count)
		return newItem
	}

	if p.initial != nil {
		// create initial number of items
		items := []interface{}{}

		// create new items
		for i := 0; i < *p.initial; i++ {
			items = append(items, p.getItem())
		}
		// return new items
		for j := len(items) - 1; j >= 0; j-- {
			p.returnItem(items[j])
		}
		p.initial = nil
	}
}

func (p *Pool) getItem() interface{} {
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
		if p.semMax != nil {
			p.semMax.Release(1)
		}
		p.syncPool.Put(wrap.Item)
	}
	itemWrapPool.Put(wrap)
}

// GetItem obtains an item from the pool.
// If the Max option is set, then this function
// will block until an item is returned back into the pool.
//
// After the item is no longer required, you must call
// either Close or MarkAsInvalid on the item.
func (p *Pool) GetItem() *ItemWrap {
	return p.getItem().(*ItemWrap)
}

// ReturnItem returns an item back to the pool.
// Usually this function is never called as the recommened
// approach is to call either Close or MarkAsInvalid on the item.
func (p *Pool) ReturnItem(x *ItemWrap) {
	p.returnItem(x)
}

// Count returns approximately the number of items in the pool (idle and in-use).
// If you want an accurate number, call runtime.GC() twice (not recommended).
func (p *Pool) Count() int {
	return int(atomic.LoadUint32(&p.count))
}