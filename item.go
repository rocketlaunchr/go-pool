package pool

// ItemWrap wraps the item returned by the pool's factory.
type ItemWrap[T any] interface {
	// Return returns the item back to the pool.
	Return()
	
	// MarkAsInvalid marks the item as invalid (eg. unusable, unstable or broken) so
	// that after it gets returned to the pool, it is discarded. It will eventually
	// get garbage collected.
	MarkAsInvalid()
	
	// Item represents the unwrapped item borrowed from the pool.
	Item() T
}

// itemWrap wraps the item returned by the pool's factory.
type itemWrap[T any] struct {
	item T

	invalid bool
	pool    interface{ returnItem(*itemWrap[T]) }
}

// Item represents the unwrapped item borrowed from the pool.
func (iw *itemWrap[T]) Item() T {
	return iw.item
}

// Return returns the item back to the pool.
func (iw *itemWrap[T]) Return() {
	iw.pool.returnItem(iw)
}

// reset restores iw to the zero value.
func (iw *itemWrap[T]) reset() {
	iw.item = *new(T)
	iw.invalid = false
	iw.pool = nil
}

// MarkAsInvalid marks the item as invalid (eg. unusable, unstable or broken) so
// that after it gets returned to the pool, it is discarded. It will eventually
// get garbage collected.
func (iw *itemWrap[T]) MarkAsInvalid() {
	iw.invalid = true
}
