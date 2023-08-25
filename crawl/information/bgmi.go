package information

import (
	"encoding/json"
	"fmt"
	CR "github.com/NullpointerW/anicat/crawl"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"net/http"
	"strings"
)

var endpoint = "search/subject/%s?type=2&start=%d&max_results=%d"

func BgmiApiSearch(searchstr string) (sid int, err error) {
	searchstr = CR.UrlEncode(searchstr)
	ed := fmt.Sprintf(endpoint, searchstr, 0,
		10)
	log.Info(log.NewUrlStruct(CR.BgmiRoot+ed), "request bgmTV search api")
	req, err := http.NewRequest("GET", CR.BgmiRoot+ed, nil)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrBgmTVApiPrefix, err)
	}
	resp, err := CR.BgmiRequest(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrBgmTVApiPrefix, err)
	}
	bsis := struct {
		List []CR.BgmiSubjIntro `json:"list"`
	}{}
	jde := json.NewDecoder(resp.Body)
	err = jde.Decode(&bsis)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrBgmTVApiPrefix, err)
	}
	if len(bsis.List) == 0 {
		return 0, fmt.Errorf("%w: %w", errs.ErrBgmTVApiPrefix, errs.ErrCrawlNotFound)
	}
	var tatget *CR.BgmiSubjIntro
	for _, bsi := range bsis.List {
		log.Debug(log.Struct{"name", bsi.NameCN}, "traverse bgmTV  matching items")
		if bsi.NameCN == searchstr {
			log.Infof(log.Struct{"matched", bsi}, "%s: matching item found", errs.ErrBgmTVApiPrefix)
			tatget = &bsi
			break
		}
	}
	if tatget == nil {
		tatget = &bsis.List[0]
	}
	return tatget.Id, err
}

func BgmTVInfoScrape(sid int) (tips map[string]string, err error) {
	url := infoBaseUrl + fmt.Sprintf(InfoAPIs["subject"], sid)
	tips, err = DoScrape(url)
	return
}

func Scrape(searchstr string) (tips map[string]string, err error) {
	sid, err := BgmiApiSearch(searchstr)
	// p, err := InfoPageScrape(searchstr)
	if err != nil {
		return tips, err
	}
	tips, err = BgmTVInfoScrape(sid)
	// url := infoBaseUrl + p
	// tips, err = DoScrape(url)
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
			// fetch origin name
			a := htmlquery.FindOne(doc, OriginNameXpath)
			tips[SubjOriginName] = htmlquery.InnerText(a)
		} else {
			err = fmt.Errorf("%w: bgmi info", errs.ErrCrawlNotFound)
			return
		}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info(log.NewUrlStruct(r.URL), "searching info from bgmTV")
	})
	c.OnError(func(_ *colly.Response, e error) {
		e = fmt.Errorf("%s: search info failed: %w", errs.ErrBgmTVApiPrefix, e)
		err = e
		log.Error(nil, err)
	})
	c.Visit(url)
	return tips, err
}

// Deprecated: use bgmi search api
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
		} else {
			err = errs.ErrCrawlNotFound
			return
		}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info(log.NewUrlStruct(r.URL), "fetching info from bgmTV")
	})
	c.OnError(func(_ *colly.Response, e error) {
		e = fmt.Errorf("fetch info from bgmTV failed: %w", e)
		err = e
		log.Error(nil, err)
	})
	c.Visit(BuildInfoSearching(CR.UrlEncode(searchstr)))
	return p, err
}

func BuildInfoSearching(s string) string {
	return fmt.Sprintf(infoBaseUrl+InfoAPIs["search"], s)
}
