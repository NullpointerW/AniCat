package crwal

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

const SearchUrl = `https://bgm.tv/subject_search/%s?cat=2`




func InfoScrape(searchstr string) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		a := htmlquery.FindOne(doc, `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@id='main'][2]/div[@class='columns clearit']/div[@id='columnSearchB']/ul[@id='browserItemList']/li[1]/div[@class='inner']/h3/a[@class='l']/@href`)
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

	if p, err := proxy.RoundRobinProxySwitcher(
		"http://127.0.0.1:7890",
	); err == nil {
		c.SetProxyFunc(p)
	}

	c.Visit(BuildInfoSearching(ConstructSearch(searchstr)))

}

func BuildInfoSearching(s string)string{
	return fmt.Sprintf(SearchUrl,s)
}