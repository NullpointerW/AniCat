package crawl

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

const (
	infoBaseUrl  = `https://bgm.tv`
	infoXpathExp = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@id='main'][2]/div[@class='columns clearit']/div[@id='columnSearchB']/ul[@id='browserItemList']/li[1]/div[@class='inner']/h3/a[@class='l']/@href`
)

var InfoAPIs = map[string]string{
	"search":  "/subject_search/%s?cat=2",
	"subject": "/subject/%d",
}

func InfoPageScrape(searchstr string) (p string, err error) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			log.Fatal(err)
		}
		a := htmlquery.FindOne(doc, infoXpathExp)
		if a != nil {
			p = htmlquery.InnerText(a)
			fmt.Println(htmlquery.InnerText(a))
		} else {
			fmt.Println("NOT FOUND")
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		log.Println("Something went wrong:", e)
		err = e
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
	return p, err

}

func BuildInfoSearching(s string) string {
	return fmt.Sprintf(infoBaseUrl+InfoAPIs["search"], s)
}
