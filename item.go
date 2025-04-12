package pool

// ItemWrap wraps the item returned by the pool's factory.
type ItemWrap[T any] struct {
	Item T

	invalid bool
	pool    interface{ returnItem(*ItemWrap[T]) }
}

// Return returns the item back to the pool.
func (iw *ItemWrap[T]) Return() {
	iw.pool.returnItem(iw)
}

// Reset restores iw to the zero value.
func (iw *ItemWrap[T]) Reset() {
	iw.Item = *new(T)
	iw.invalid = false
	iw.pool = nil
}

// MarkAsInvalid marks the item as invalid (eg. unusable, unstable or broken) so
// that after it gets put back in the pool, it is discarded. It will eventually
// get garbage collected.
func (iw *ItemWrap[T]) MarkAsInvalid() {
	iw.invalid = true
}
