package cover

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	CR "github.com/NullpointerW/mikanani/crawl"
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
	c.OnResponse(func(r *colly.Response) {
		exp := coverXpathExp
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
		log.Printf("cover file downloaded size:%d", wn)
		if e != nil {
			fmt.Println(e)
			err = e
		}
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
	c := CR.NewCollector()
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

	parseParam := CR.ConstructSearch(coverName)
	c.Visit(fmt.Sprintf(coverSearchUrl, parseParam))

	return
}
