package crwal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

type CoverScraper interface {
	Scrap(filePath, CoverName string)
}

const CoverSearchUrl = `https://movie.douban.com/j/subject_suggest?q=%s`

func TouchCoverImg(fpath, cover string) {
	u := coverImgScrap(cover)
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		exp := `/html/body/div[@id='wrapper']/div[@id='content']/div[@class='grid-16-8 clearfix']/div[@class='article']/ul[@class='poster-col3 clearfix']/li[1]/div[@class='cover']/a/img/@src`
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		a := htmlquery.FindOne(doc, exp)
		m := htmlquery.InnerText(a)
		dl := strings.ReplaceAll(m, `/m/`, `/l/`)
		fmt.Println(dl)
		//download
		resp, err := http.Get(dl)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		defer f.Close()
		wn, err := io.Copy(f, resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(wn)
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

	c.Visit(u)
}

func coverImgScrap(coverName string) (cUrl string) {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		jsonstr := string(r.Body)
		subjUrl := gjson.Get(jsonstr, "0").Get("url").String()
		u, _ := url.Parse(subjUrl)
		u.RawQuery = ""
		cUrl = u.String() + `photos?type=R`
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
	parseParam := ConstructSearch(coverName)
	c.Visit(fmt.Sprintf(CoverSearchUrl, parseParam))
	return
}

// func ConstructCoverSearch(s string) string {
// 	p := strings.ReplaceAll(s, " ", "+")
// 	return fmt.Sprintf(CoverSearchUrl, p)
// }
