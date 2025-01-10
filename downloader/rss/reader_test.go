package rss

import (
	"fmt"
	"strings"
	"testing"
)

func TestReaderOnce(t *testing.T) {
	r := NewReader("https://mikanani.me/RSS/Bangumi?bangumiId=3375&subgroupid=34", nil, func(n string) bool {
		return strings.Contains(n, "1080p")
	})
	it, ok, err := r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

	it, ok, err = r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

	it, ok, err = r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

}

func TestReaderOnce2(t *testing.T) {
	r := NewReader("https://mikanani.me/RSS/Bangumi?bangumiId=3375&subgroupid=34", nil, nil)
	it, ok, err := r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

	it, ok, err = r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

	it, ok, err = r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

}
func TestRead(t *testing.T) {
	r := NewReader("https://mikanani.me/RSS/Bangumi?bangumiId=3375&subgroupid=34", nil, nil)
	its, ok, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(its)
	its, ok, err = r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(its)

	it, ok, err := r.ReadOne()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("read eof")
	}
	fmt.Println(it)

	fmt.Println(r.Guids())

}
