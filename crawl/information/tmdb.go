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

var (
	root         = `/html/body[@class='zh v4']/div[1]/main[1]/section[@class='main_content search_results']/div[@class='column_wrapper reverse']/div[@class='content_wrapper']/div[@class='white_column']/section[@class='panel']/div[1]/div[@class='results flex']/div[1]/div[@class='wrapper']/div[@class='details']/div[@class='wrapper']/div[@class='title']`
	nameXpathExp = root + `/div/a[@class='result']/h2`
	dateXpathExp = root + `/span[@class='release_date']`
)

func FloderSearch(typ, searchstr string) (name, date string, err error) {
	url := TMDB_HOST + fmt.Sprintf(TMDBAPIs["search"], TMDB_TYP_TV, "")
	if typ == TMDB_TYP_MOVIE {
		url = TMDB_HOST + fmt.Sprintf(TMDBAPIs["search"], TMDB_TYP_MOVIE, "")
	}

	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		nameH2 := htmlquery.FindOne(doc, nameXpathExp)
		if nameH2 != nil {
			name = htmlquery.InnerText(nameH2)
		} else {
			err = fmt.Errorf("%w:TMDB info not found,search str:%s", errs.ErrCrawlNotFound, searchstr)
			return
		}
		dateSpan := htmlquery.FindOne(doc, dateXpathExp)
		if dateSpan != nil {
			date = htmlquery.InnerText(dateSpan)
		} else {
			err = fmt.Errorf("%w:TMDB info not found", errs.ErrCrawlNotFound)
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
		r.Headers.Set("User-Agent", agent)
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		log.Println("Something went wrong:", e)
		err = e
	})

	c.Visit(url + CR.ConstructSearch(searchstr))
	return name, date, err
}
