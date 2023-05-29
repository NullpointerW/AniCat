package util

import (
	"fmt"
	"testing"
)

type tests struct{}

func (t tests) Print() {

}

func TestXxx(t *testing.T) {
	ti, err := ParseTime("2009年4月12日")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(ti)
}

func TestTuple(t *testing.T) {
	tuple := NewTuple("a", 3)
	// tuple2 := NewTuple(4, 3)
	f := func(t Tuple[string, int]) {

	}
	f(tuple)
	// f(tuple2)

	s := tuple.Get0()
	i := tuple.Get1()
	fmt.Println(s, i)
	for i := 0; i < 20; i++ {
		blank := &struct{}{}
		tp := &tests{}
		fmt.Println(tp == (*tests)(blank))
		fmt.Printf("%p \n", tp)
		fmt.Printf("%p", blank)
	}
}
