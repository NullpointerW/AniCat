package crwal

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

const CoverSearchUrl = `https://search.douban.com/movie/subject_search?search_text=%s&cat=1002`

func CoverImgScrap(coverName string) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		a := htmlquery.FindOne(doc, `/html/body/div[@id='wrapper']/div[@id='root']/div[@class='sc-gqjmRU gsUvmP']/div[@class='_xbuf0ntj7']/div[@class='_9fled6mja']/div[1]/div[@class='sc-bZQynM jReLEq sc-bxivhb eJWSlY'][1]/div[@class='item-root']/div[@class='detail']/div[@class='title']/a/@href`)
		if a != nil {
			fmt.Println(htmlquery.InnerText(a))
		} else {
			fmt.Println("NOT FOUND")
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// c.OnScraped(func(r *colly.Response) {
	//     fmt.Printf("Finished total count:%d /n", count)
	// })

	c.Visit(BuildSearching(ConstructCoverSearch(coverName)))
}

func ConstructCoverSearch(s string) string {
	p := strings.ReplaceAll(s, " ", "+")
	return fmt.Sprintf(CoverSearchUrl, p)
}
