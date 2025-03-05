package util

import (
	"slices"
	// "sync"
	"sync/atomic"
)

func SetSubtract[T comparable, V any](a map[T]V, b map[T]struct{}) (res map[T]V) {
	res = make(map[T]V)
	for e, v := range a {
		if _, ex := b[e]; !ex {
			res[e] = v
		}
	}
	return
}
func SliceDelete[S ~[]E, E any](s S, i int) S {
	return slices.Delete(s, i, i+1)
}

// ListView is a thread-safe list view that allows only a single goroutine to append elements
// while permitting multiple goroutines to read concurrently.
// This data structure is particularly useful when sharing a list across multiple goroutines
// where only one goroutine is responsible for adding new elements.
type ListView[T any] struct {
	ptr atomic.Pointer[[]T]
}

// Append adds a new element to the list.
// Note: This method should only be called by a single goroutine as it's not designed
// for concurrent append operations.
func (lv *ListView[T]) Append(val T) {
	var ll *[]T
	if ll = lv.ptr.Load(); ll == nil {
		new := make([]T, 0, 1)
		ll = &new
	}
	new := append(*ll, val)
	lv.ptr.Store(&new)
}

// List returns a copy of all elements in the list.
// This method is safe to be called concurrently by multiple goroutines.
func (lv *ListView[T]) List() []T {
	ll := lv.ptr.Load()
	if ll == nil {
		return ([]T)(nil)
	}
	return *ll
}

func NewListView[T any](ls []T) *ListView[T] {
	lv := &ListView[T]{}
	lv.ptr.Store(&ls)
	return lv
}
