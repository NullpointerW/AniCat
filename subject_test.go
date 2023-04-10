package main

import (
	"fmt"
	"os"
	"testing"
)

func TestOS(t *testing.T) {
	_,err:=os.ReadDir("rows")
	fmt.Println(err)
}