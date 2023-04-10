package crwal

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func TestCrwal(t *testing.T) {
	Scrape("凉宫春日")
	Scrape("lycoris Recoil")
}

func TestCoverCrwal(t *testing.T) {
	fmt.Println(coverImgScrape("凉宫 春日"))
	// https://www.douban.com/search?q=%E5%87%89%E5%AE%AB%E6%98%A5%E6%97%A5

	// /html[@class='ua-windows ua-webkit']/body/div[@id='wrapper']/div[@id='content']/div[@class='grid-16-8 clearfix']/div[@class='article']/div[@class='search-result']/div[@class='result-list'][1]/div[@class='result'][1]/div[@class='content']/div[@class='title']/h3/a/@href
}

func TestJsonArray(t *testing.T) {
	jsonstr := `[1,2,3]`
	fmt.Println(gjson.Get(jsonstr, "0").Int())
	// https://movie.douban.com/subject/4074292/?suggest=%E5%87%89%E5%AE%AB+%E6%98%A5%E6%97%A5
}

func TestTouchCoverImg(t *testing.T) {
	err := DOUBANCoverScraper.Scrape("cover_testing/cover7.jpg", "奇蛋物语")
	fmt.Println(err)
}
