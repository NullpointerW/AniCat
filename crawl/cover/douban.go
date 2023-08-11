package cover

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	CR "github.com/NullpointerW/anicat/crawl"
	util "github.com/NullpointerW/anicat/utils"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func TouchCoverImg(fpath, cover string) (err error) {
	u, err := coverImgScrape(cover)
	if err != nil {
		return err
	}
	c := CR.NewCollector()
	c.SetRequestTimeout(5 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"

		r.Headers.Set("User-Agent", agent)
		r.Headers.Set("Sec-Ch-Ua", `"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`)
		r.Headers.Set("Sec-Ch-Ua-Platform", `"Android"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?1")

		util.Debugf("%#+v", r.Headers)
	})
	c.OnResponse(func(r *colly.Response) {
		exp := DouBancoverXpathExp
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		a := htmlquery.FindOne(doc, exp)
		m := htmlquery.InnerText(a)
		dl := strings.ReplaceAll(m, `/m/`, `/l/`)
		log.Println("DOUBAN cover: file url:", dl)
		//download

		resp, e := http.Get(dl)
		if e != nil {
			err = e
			return
		}

		e = CR.Downloadfile(fpath, resp.Body)
		if e != nil {
			err = e
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("scraping cover from DOUBAN", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		e = fmt.Errorf("scrap DOUBAN cover failed: %w", e)
		err = e
		log.Println(e)
	})

	c.Visit(u)

	return err
}

func coverImgScrape(coverName string) (cUrl string, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		jsonstr := string(r.Body)
		subjUrl := gjson.Get(jsonstr, "0").Get("url").String()
		u, _ := url.Parse(subjUrl)
		u.RawQuery = ""
		cUrl = u.String() + `photos?type=R`
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("searching cover from DOUBAN", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		e = fmt.Errorf("search cover failed: %w", e)
		err = e
		log.Println(err)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Printf("coverScrapUrl=%s \n", cUrl)
	})

	parseParam := CR.UrlEncode(coverName)
	c.Visit(fmt.Sprintf(DouBancoverSearchUrl, parseParam))

	return
}
