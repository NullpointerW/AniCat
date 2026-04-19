package information

import (
	"fmt"
	CR "github.com/NullpointerW/anicat/crawl"
	sel "github.com/NullpointerW/anicat/crawl/selector"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"strings"
)

const (
	nameXpathExpDefault = `//h2`
	dateXpathExpDefault = `//span[contains(@class,'release_date')]`
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
		pageURL := url + CR.UrlEncode(searchstr)
		nameH2 := htmlquery.FindOne(doc, tmdbXpath("title", nameXpathExpDefault))
		if nameH2 != nil {
			name = htmlquery.InnerText(nameH2)
		} else {
			err = fmt.Errorf("%w: TMDB info not found,search str=%s", errs.ErrCrawlNotFound, searchstr)
			sel.TriggerHeal("tmdb", "title", pageURL)
			return
		}
		dateSpan := htmlquery.FindOne(doc, tmdbXpath("release_date", dateXpathExpDefault))
		if dateSpan != nil {
			date = htmlquery.InnerText(dateSpan)
		} else {
			err = fmt.Errorf("%w: TMDB info not found", errs.ErrCrawlNotFound)
			sel.TriggerHeal("tmdb", "release_date", pageURL)
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
		r.Headers.Set("User-Agent", agent)
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
		log.Debug(log.NewUrlStruct(r.URL, "source", "TMDB", "searchStr", searchstr), "fetching folderInfo")
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = fmt.Errorf("search/fetch folder info from TMDB failed: %w", e)
		log.Error(log.Struct{"url", url, "searchStr", searchstr}, err)
	})

	c.Visit(url + CR.UrlEncode(searchstr))
	return name, date, err
}
