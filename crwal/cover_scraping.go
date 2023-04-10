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
	Scrape(filePath, CoverName string)
}

const CoverSearchUrl = `https://movie.douban.com/j/subject_suggest?q=%s`

func TouchCoverImg(fpath, cover string) (err error) {
	u, err := coverImgScrape(cover)
	if err != nil {
		return err
	}
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	c.OnResponse(func(r *colly.Response) {
		exp := `/html/body/div[@id='wrapper']/div[@id='content']/div[@class='grid-16-8 clearfix']/div[@class='article']/ul[@class='poster-col3 clearfix']/li[1]/div[@class='cover']/a/img/@src`
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			log.Fatal(e)
		}
		a := htmlquery.FindOne(doc, exp)
		m := htmlquery.InnerText(a)
		dl := strings.ReplaceAll(m, `/m/`, `/l/`)
		fmt.Println(dl)
		//download
		resp, e := http.Get(dl)
		if e != nil {
			log.Fatal(e)
			err = e
		}
		f, e := os.Create(fpath)
		if e != nil {
			log.Fatal(e)
			err = e
		}
		defer resp.Body.Close()
		defer f.Close()
		wn, e := io.Copy(f, resp.Body)
		if e != nil {
			fmt.Println(e)
			err = e
		}
		fmt.Println(wn)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit(u)
	return err
}

func coverImgScrape(coverName string) (cUrl string, err error) {
	c := colly.NewCollector()
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

	c.OnError(func(_ *colly.Response, e error) {
		log.Println("Something went wrong:", e)
		err = e
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("coverScrapUrl=%s \n", cUrl)
	})

	parseParam := ConstructSearch(coverName)
	c.Visit(fmt.Sprintf(CoverSearchUrl, parseParam))
	return
}
