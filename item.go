package pool

import "sync"

// itemWrapPool is a global pool of *ItemWrap structs.
var itemWrapPool = sync.Pool{New: func() interface{} { return new(ItemWrap) }}

// ItemWrap wraps the item returned by the pool's factory.
type ItemWrap struct {
	Item interface{}

	invalid bool
	pool    *Pool
}

// Return returns the item back to the pool.
func (iw *ItemWrap) Return() {
	iw.pool.returnItem(iw)
}

// MarkAsInvalid marks the item as invalid (eg. unusable, unstable or broken) so
// that after it gets put back in the pool it is discarded. It will eventually
// get garbage collected.
func (iw *ItemWrap) MarkAsInvalid() {
	iw.invalid = true
}
