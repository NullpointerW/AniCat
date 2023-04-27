package resource

import (
	"fmt"
	"log"
	"strings"

	CR "github.com/NullpointerW/mikanani/crawl"
	"github.com/NullpointerW/mikanani/errs"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

func Scrape(searchstr string) (url string, isrss bool, err error) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			err = e
			return
		}
		a := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/ul[@class='list-inline an-ul']/li/a//@href`)
		if a != nil {
			ep, e := scrapeRssEndPoint(htmlquery.InnerText(a))
			if e != nil {
				err = e
				return
			}
			url = resourcesBaseUrl + ep
			isrss = true
		} else {
			log.Println("RSS_NOT_FOUND")
			fn := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]/a[@class='magnet-link-wrap']`)
			size := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[2]`)
			t := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[3]`)
			torr := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[4]/a/@href`)
			mglink := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]/a[2]/@data-clipboard-text`)
			log.Printf("file_name:%s,size:%s,update_time=%s,torrent:%s,magnetLink=%s",
				htmlquery.InnerText(fn),
				htmlquery.InnerText(size),
				htmlquery.InnerText(t),
				htmlquery.InnerText(torr),
				htmlquery.InnerText(mglink),
			)
			url = htmlquery.InnerText(mglink)
			isrss = false
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	CR.SetProxy(c)

	c.Visit(BuildSearching(CR.ConstructSearch(searchstr)))

	return
}

func scrapeRssEndPoint(endpoint string) (rssep string, err error) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
		}
		xpathExp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/div[@class='subgroup-text'][1]/a[@class='mikan-rss']/@href`
		a := htmlquery.FindOne(doc, xpathExp)
		if a == nil {
			err = errs.ErrCrawlNotFound
			return
		} else {
			rssep = htmlquery.InnerText(a)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	CR.SetProxy(c)

	c.Visit(resourcesBaseUrl + endpoint)

	return
}

func BuildSearching(s string) string {
	return resourcesBaseUrl + ResourceAPIs["search"] + s
}
