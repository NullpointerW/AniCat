package test

import (
	"fmt"
	"testing"

	C "github.com/NullpointerW/anicat/crawl/cover"
)

func TestTouchCoverImg(t *testing.T) {
	err := C.DOUBANCoverScraper.Scrape("cover.jpg", "雪之少女")
	fmt.Println(err)
}
