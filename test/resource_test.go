package test

import (
	R "github.com/NullpointerW/mikanani/crawl/resource"
	"testing"
)

func TestCrwal(t *testing.T) {
	R.Scrape("凉宫春日")
	R.Scrape("lycoris Recoil")
}
