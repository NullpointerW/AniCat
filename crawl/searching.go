package crawl

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

const resourcesBaseUrl = `https://mikanani.me`

var ResourceAPIs = map[string]string{
	"search": "/Home/Search?searchstr=",
}

func Scrape(searchstr string) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		a := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/ul[@class='list-inline an-ul']/li/a//@href`)
		if a != nil {
			fmt.Println(htmlquery.InnerText(a))
		} else {
			fmt.Println("NOT FOUND")
			fn := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]/a[@class='magnet-link-wrap']`)
			size := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[2]`)
			t := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[3]`)
			torr := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[4]/a/@href`)
			fmt.Printf("file_name:%s,size:%s,update_time=%s,torrent:%s",
				htmlquery.InnerText(fn),
				htmlquery.InnerText(size),
				htmlquery.InnerText(t),
				htmlquery.InnerText(torr))
		}
		// for _, node := range divNodes {
		// 	url := htmlquery.FindOne(node, "./h1[@class='post-title entry-title']/a/@href")
		// 	text := htmlquery.FindOne(node, "./h1[@class='post-title entry-title']/a")

		// 	fmt.Println(htmlquery.InnerText(url))
		// 	fmt.Println(htmlquery.InnerText(text))
		// 	rows[htmlquery.InnerText(text)] = htmlquery.InnerText(url)
		// 	count++
		// }
		// preFound := htmlquery.Find(doc, `/html/body/div[@id='outer-wrapper']/div[@id='wrap2']/div[@id='content-wrapper']/div[@id='main-wrapper']/div[@id='main']/div[@id='Blog1']/div[@id='blog-pager']/span[@id='blog-pager-older-link']/a[@id='Blog1_blog-pager-older-link']/@href`)
		// if len(preFound) != 0 {
		// 	pre := preFound[0]
		// 	link := htmlquery.InnerText(pre)
		// 	c.Visit(link)
		// }
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

	c.Visit(BuildSearching(ConstructSearch(searchstr)))

}

func ConstructSearch(s string) (utoa string) {
	a := url.QueryEscape(strings.ReplaceAll(s, " ", "+"))
	utoa = strings.ReplaceAll(a, "%2B", "+")
	return
}

func BuildSearching(s string) string {
	return resourcesBaseUrl + ResourceAPIs["search"] + s
}
