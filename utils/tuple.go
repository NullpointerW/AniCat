package util

type Tuple[F, S any] struct {
	slot1 F
	slot2 S
}

func (t *Tuple[F, S]) Get0() F {
	return t.slot1
}

func (t *Tuple[F, S]) Get1() S {
	return t.slot2
}

func NewTuple[F, S any](a1 F, a2 S) Tuple[F, S] {
	return Tuple[F, S]{
		a1,
		a2,
	}
}
