package resource

import (
	"fmt"
	CR "github.com/NullpointerW/mikanani/crawl"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

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
			mglink := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]/a[2]/@data-clipboard-text`)
			fmt.Printf("file_name:%s,size:%s,update_time=%s,torrent:%s,magnetLink=%s",
				htmlquery.InnerText(fn),
				htmlquery.InnerText(size),
				htmlquery.InnerText(t),
				htmlquery.InnerText(torr),
				htmlquery.InnerText(mglink),
			)
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

}

func BuildSearching(s string) string {
	return resourcesBaseUrl + ResourceAPIs["search"] + s
}
