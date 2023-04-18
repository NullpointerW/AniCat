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
	infoBaseUrl      = `https://bgm.tv`
	infoPageXpathExp = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@id='main'][2]/div[@class='columns clearit']/div[@id='columnSearchB']/ul[@id='browserItemList']/li[1]/div[@class='inner']/h3/a[@class='l']/@href`
	infoXpathExp     = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@class='mainWrapper']/div[@class='columns clearit']/div[@id='columnSubjectHomeA']/div[@id='bangumiInfo']/div[@class='infobox']/ul[@id='infobox']/li`
)

var InfoAPIs = map[string]string{
	"search":  "/subject_search/%s?cat=2",
	"subject": "/subject/%d",
}

func InfoScrape(searchstr string) (tips map[string]string, err error) {
	tips = make(map[string]string)
	p,err:=InfoPageScrape(searchstr)
	if err!=nil{
		return tips,err
	}
	url:=infoBaseUrl+p
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			log.Fatal(err)
		}
		ls := htmlquery.Find(doc, infoXpathExp)

		if ls != nil {
			for _, l := range ls {
				t := htmlquery.FindOne(l, "./span")
				tt := htmlquery.InnerText(t)
				c := htmlquery.InnerText(l)
				tips[tt] = c
			}
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

	c.Visit(url)
	return tips, err
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
		a := htmlquery.FindOne(doc, infoPageXpathExp)
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
