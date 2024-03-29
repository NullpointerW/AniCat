package util

import (
	"fmt"
	"testing"
)

type tests struct{}

func (t tests) Print() {

}

func TestXxx(t *testing.T) {
	ti, err := ParseTime("2009年4月12日", YMDParseLayout)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(ti)
}

func TestParseShortTime(t *testing.T) {
	sstr, err := ParseShortTime("2023年1月9日")
	if err != nil {
		t.Error(t)
		t.FailNow()
	}
	fmt.Println(sstr)
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

func TestIsNumber(t *testing.T) {
	t.Log(IsNumber("十二"))
}

func TestConvertZHCN(t *testing.T) {
	t.Log(ConvertZhCnNumbToa("三十五"))
}

func TestCheckZH(t *testing.T) {
	t.Log(CheckZhCn("av2"))
}

func TestTrimGetEpi(t *testing.T) {
	t.Log(TrimExtensionAndGetEpi("天国大魔镜 S01E02.mp4"))
}

func TestTrimGetEpi2(t *testing.T) {
	// fmt.Printf("%02d", 1)
	fmt.Printf("%02s", "02")
}
