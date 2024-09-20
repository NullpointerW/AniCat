package util

import "slices"

func SetSubtract[T comparable,V any](a map[T]V, b map[T]struct{}) (res map[T]V) {
	res = make(map[T]V)
	for e,v := range a {
		if _, ex := b[e]; !ex {
			res[e] =v
		}
	}
	return
}
func SliceDelete[S ~[]E, E any](s S, i int) S {
	return slices.Delete(s, i, i+1)
}
