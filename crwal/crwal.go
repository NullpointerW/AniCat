package crwal
import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func Scrape()(rows map[string]string) {
	rows=make(map[string]string)
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	var count int
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		divNodes := htmlquery.Find(doc, `/html/body/div[@id='outer-wrapper']/div[@id='wrap2']/div[@id='content-wrapper']/div[@id='main-wrapper']/div[@id='main']/div[@id='Blog1']/div[@class='blog-posts hfeed']/div[@class='post hentry uncustomized-post-template']`)
		for _, node := range divNodes {
			url := htmlquery.FindOne(node, "./h1[@class='post-title entry-title']/a/@href")
			text := htmlquery.FindOne(node, "./h1[@class='post-title entry-title']/a")
			
			fmt.Println(htmlquery.InnerText(url))
			fmt.Println(htmlquery.InnerText(text))
			rows[htmlquery.InnerText(text)]=htmlquery.InnerText(url)
			count++
		}
		preFound := htmlquery.Find(doc, `/html/body/div[@id='outer-wrapper']/div[@id='wrap2']/div[@id='content-wrapper']/div[@id='main-wrapper']/div[@id='main']/div[@id='Blog1']/div[@id='blog-pager']/span[@id='blog-pager-older-link']/a[@id='Blog1_blog-pager-older-link']/@href`)
		if len(preFound) != 0 {
			pre := preFound[0]
			link := htmlquery.InnerText(pre)
			c.Visit(link)
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

	c.Visit("https://program-think.blogspot.com/")

	fmt.Println(count)
	return rows

}
