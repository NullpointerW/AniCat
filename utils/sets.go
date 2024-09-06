package util

func SetSubtract[T comparable](a map[T]struct{}, b map[T]struct{}) (res map[T]struct{}) {
	res = make(map[T]struct{})
	for e := range a {
		if _, ex := b[e]; !ex {
			res[e] = struct{}{}
		}
	}
	return
}
