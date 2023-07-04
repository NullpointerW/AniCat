package information

import (
	"fmt"
	"log"
	"strings"

	CR "github.com/NullpointerW/anicat/crawl"
	"github.com/NullpointerW/anicat/errs"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

func BgmTVInfoScrape(sid int) (tips map[string]string, err error) {
	url := infoBaseUrl + fmt.Sprintf(InfoAPIs["subject"], sid)
	tips, err = DoScrape(url)
	return
}

func Scrape(searchstr string) (tips map[string]string, err error) {
	p, err := InfoPageScrape(searchstr)
	if err != nil {
		return tips, err
	}
	url := infoBaseUrl + p
	tips, err = DoScrape(url)
	return
}

func DoScrape(url string) (tips map[string]string, err error) {
	tips = make(map[string]string)
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		ls := htmlquery.Find(doc, infoXpathExp)

		if ls != nil {
			for _, l := range ls {
				t := htmlquery.FindOne(l, "./span")
				tt := htmlquery.InnerText(t)
				lt := htmlquery.InnerText(l)
				lt = strings.Replace(lt, tt, "", 1)
				tt = strings.TrimSuffix(tt, ": ")
				pre, e := tips[tt]
				if e {
					lt += "|" + pre
				}
				tips[tt] = lt
			}
			s := strings.Split(url, `/`)
			sid := s[len(s)-1]
			tips["sid"] = sid
		} else {
			log.Println("nofound bgmi cover")
			err = errs.ErrCrawlNotFound
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

	CR.SetProxy(c)

	c.Visit(url)
	return tips, err
}

func InfoPageScrape(searchstr string) (p string, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		a := htmlquery.FindOne(doc, infoPageXpathExp)
		if a != nil {
			p = htmlquery.InnerText(a)
			fmt.Println(htmlquery.InnerText(a))
		} else {
			// fmt.Println("NOT FOUND")
			err = errs.ErrCrawlNotFound
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

	c.Visit(BuildInfoSearching(CR.ConstructSearch(searchstr)))
	return p, err
}

func BuildInfoSearching(s string) string {
	return fmt.Sprintf(infoBaseUrl+InfoAPIs["search"], s)
}
