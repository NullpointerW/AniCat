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

func TestMap(t *testing.T){
	var m map[int]int =map[int]int{}
	m1:=m
	m[1]=1
	m[2]=2
	fmt.Println(len(m1))
	fmt.Println(len(m))
}