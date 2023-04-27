package util

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	ti, err := ParseTime("2009年4月12日")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(ti)
}
