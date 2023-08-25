package resource

import (
	"fmt"
	CR "github.com/NullpointerW/anicat/crawl"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"strings"
)

type Option struct {
	Index int
	Group string
}

type RssGroup struct {
	Name  string
	Items []Item
}

type Item struct {
	Name       string
	Size       string
	UpdateTime string
}

func Scrape(searchstr string, opt Option) (url, bgmUrl string, isrss bool, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			err = e
			return
		}
		a := htmlquery.Find(doc, MikanRssLiXpath)
		// command is add ... -i ,even if rss source has been found,show the search list
		if opt.Index > 0 {
			a = nil
		}
		if len(a) != 0 {
			ep, bgmurl, e := scrapeRssEndPoint(selectRss(a, searchstr), opt)
			if e != nil {
				err = e
				return
			}
			url = resourcesBaseUrl + ep
			isrss = true
			bgmUrl = bgmurl
		} else {
			log.Info(log.Struct{"search name", searchstr}, "rss resource not found")
			if opt.Index <= 0 {
				opt.Index = 1
			}
			mglinkTemp := `/html/body[@class='main']/
			div[@id='sk-container']/
			div[@class='central-container']/
			table[@class='table table-striped tbl-border fadeIn']/
			tbody/
			tr[@class='js-search-results-row'][%d]
			/td[1]/a[2]/@data-clipboard-text`
			mglink := htmlquery.FindOne(doc, fmt.Sprintf(mglinkTemp, opt.Index))
			if mglink != nil {
				url = htmlquery.InnerText(mglink)
				log.Info(log.Struct{"search name", searchstr, "index", opt.Index, "magnet", url}, "scraped resource")
				isrss = false
			} else {
				err = fmt.Errorf("%w: %s", errs.ErrCrawlNotFound, searchstr)
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info(log.Struct{"url", r.URL}, "searching resource from mikan")
	})
	c.OnError(func(_ *colly.Response, e error) {
		err = fmt.Errorf("search failed from mikan: %w", e)
		log.Error(nil, err)
	})
	c.Visit(BuildSearching(CR.UrlEncode(searchstr)))
	return
}

func selectRss(nodes []*html.Node, cmp string) string {
	var targer *html.Node
	for i, n := range nodes {
		div := htmlquery.FindOne(n, `/a/div/div/div`)
		if name := htmlquery.InnerText(div); name == cmp || i == 0 {
			targer = htmlquery.FindOne(n, `/a/@href`)
		}
	}
	return htmlquery.InnerText(targer)
}

func scrapeRssEndPoint(endpoint string, opt Option) (rssUrl, bgmurl string, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		defXpathExp := `/html/body[@class='main']/div[@id='sk-container']/
		div[@class='central-container']/
		div[@class='subgroup-text'][1]/
		a[@class='mikan-rss']/@href`
		if opt.Group == "" {
			a := htmlquery.FindOne(doc, defXpathExp)
			if a == nil {
				err = errs.ErrCrawlNotFound
				return
			} else {
				rssUrl = htmlquery.InnerText(a)
			}
		} else {
			tg := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/div[@class='subgroup-text']`
			ds := htmlquery.Find(doc, tg)
			var hitgrp *html.Node
			for _, d := range ds {
				grpn := htmlquery.FindOne(d, `/a[1]`)
				if grpn == nil || htmlquery.InnerText(grpn) == "" {
					grpn = htmlquery.FindOne(d, `/div[@class='dropdown']/div[@class='dropdown-toggle material-dropdown__btn']/span[1]`)
				}
				if grpn == nil {
					log.Debug(log.Struct{"node", d}, "not scrap group name")
					continue
				} else {
					actl := strings.ToLower(htmlquery.InnerText(grpn))
					expt := strings.ToLower(opt.Group)
					if expt == actl {
						hitgrp = d
						break
					}
				}
			}
			if hitgrp != nil {
				a := htmlquery.FindOne(hitgrp, `/a[@class='mikan-rss']/@href`)
				if a == nil {
					err = errs.ErrBgmUrlNotFoundOnMikan
					return
				} else {
					rssUrl = htmlquery.InnerText(a)
				}
			} else {
				a := htmlquery.FindOne(doc, defXpathExp)
				if a == nil {
					err = errs.ErrCrawlNotFound
					return
				} else {
					rssUrl = htmlquery.InnerText(a)
				}
			}
		}
		bgmXpathExp := `/html/body[@class='main']/div[@id='sk-container']/
		div[@class='pull-left leftbar-container']/
		p[@class='bangumi-info'][last()]/
		a/@href`
		a := htmlquery.FindOne(doc, bgmXpathExp)
		if a == nil {
			err = errs.ErrBgmUrlNotFoundOnMikan
			return
		} else {
			bgmurl = htmlquery.InnerText(a)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info(log.NewUrlStruct(r.URL), "found rss resource,fetching...")
	})
	c.OnError(func(_ *colly.Response, e error) {
		err = fmt.Errorf("fetch rss resource from mikan failed: %w", e)
		log.Error(nil, err)
	})
	c.Visit(resourcesBaseUrl + endpoint)
	return
}

func ListScrape(searchstr string, t LsTyp, searchls bool) (res any, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			err = e
			return
		}
		a := htmlquery.Find(doc, MikanRssLiXpath)
		if searchls {
			a = nil
		}
		if a != nil {
			res, e = scrapeRssList(selectRss(a, searchstr), t)
			if e != nil {
				err = e
				return
			}
		} else {
			log.Info(log.Struct{"search name", searchstr, "type", "torrent"}, "command lsi")
			switch t {
			case Ls:
				fnTemp := `/html/body[@class='main']/div[@id='sk-container']/
				div[@class='central-container']/
				table[@class='table table-striped tbl-border fadeIn']/
				tbody/
				tr[@class='js-search-results-row'][%d]/
				td[1]/
				a[@class='magnet-link-wrap']`
				szTemp := `/html/body[@class='main']/div[@id='sk-container']/
				div[@class='central-container']/
				table[@class='table table-striped tbl-border fadeIn']/
				tbody/
				tr[@class='js-search-results-row'][%d]/
				td[2]`
				uptTemp := `/html/body[@class='main']/div[@id='sk-container']/
				div[@class='central-container']/
				table[@class='table table-striped tbl-border fadeIn']/
				tbody/
				tr[@class='js-search-results-row'][%d]/
				td[3]`
				/*
				 torrTemp := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/
				 div[@class='central-container']/
				 table[@class='table table-striped tbl-border fadeIn']/
				 tbody/
				 tr[@class='js-search-results-row'][%d]/
				 td[4]/a/@href`)
				*/
				nodes := htmlquery.Find(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/
				table[@class='table table-striped tbl-border fadeIn']/
				tbody/tr[@class='js-search-results-row']`)
				var items []Item
				for i, _ := range nodes {
					fn := htmlquery.FindOne(doc, fmt.Sprintf(fnTemp, i+1))
					sz := htmlquery.FindOne(doc, fmt.Sprintf(szTemp, i+1))
					upt := htmlquery.FindOne(doc, fmt.Sprintf(uptTemp, i+1))
					it := Item{}
					it.Name = InnerTextSafety(fn)
					it.Size = InnerTextSafety(sz)
					it.UpdateTime = InnerTextSafety(upt)
					items = append(items, it)
				}
				if len(items) == 0 {
					err = fmt.Errorf("%w: %s %s", errs.ErrCrawlNotFound, t.String(), searchstr)
				} else {
					res, err = items, nil
				}
			case LSGroup:
				err = fmt.Errorf("%w: search item `%s` is torrent type ", errs.ErrLsGroupUnavailableOnTorr, searchstr)
			default:
				err = errs.ErrUnknownResCrawlLsType
			}

		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info(log.NewUrlStruct(r.URL), "searching list from mikan")
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = fmt.Errorf("search resource list failed: %w", e)
		log.Error(nil, err)
	})

	c.Visit(BuildSearching(CR.UrlEncode(searchstr)))

	return
}

func scrapeRssList(endpoint string, t LsTyp) (res any, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
			return
		}
		if t != Ls && t != LSGroup {
			err = errs.ErrUnknownResCrawlLsType
			return
		}
		trsTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/
		table[@class='table table-striped tbl-border fadeIn'][%d]/
		tbody/tr`
		itnExp := `/td[1]/a[@class='magnet-link-wrap']`
		szExp := `/td[2]`
		uptExp := `/td[3]`
		var rgs []RssGroup
		tg := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/div[@class='subgroup-text']`
		ds := htmlquery.Find(doc, tg)
		for i, d := range ds {
			rg := RssGroup{}
			grpn := htmlquery.FindOne(d, `/a[1]`)
			if grpn == nil || htmlquery.InnerText(grpn) == "" {
				grpn = htmlquery.FindOne(d, `/div[@class='dropdown']/div[@class='dropdown-toggle material-dropdown__btn']/span[1]`)
			}
			if grpn == nil {
				rg.Name = strings.Fields(htmlquery.InnerText(d))[0]
			} else {
				rg.Name = htmlquery.InnerText(grpn)
			}
			trs := htmlquery.Find(doc, fmt.Sprintf(trsTemp, i+1))
			var its []Item
			for _, tr := range trs {
				itname := InnerTextSafety(htmlquery.FindOne(tr, itnExp))
				sz := InnerTextSafety(htmlquery.FindOne(tr, szExp))
				upt := InnerTextSafety(htmlquery.FindOne(tr, uptExp))
				item := Item{}
				item.Name = itname
				item.Size = sz
				item.UpdateTime = upt
				its = append(its, item)
			}
			rg.Items = its
			rgs = append(rgs, rg)
		}
		if len(rgs) == 0 {
			err = fmt.Errorf("%w: cannot found any rss group from %s ", errs.ErrCrawlNotFound, r.Request.URL.String())
			return
		}
		if t == LSGroup {
			var rgnls []string
			for _, rg := range rgs {
				rgnls = append(rgnls, rg.Name)
			}
			res = rgnls
			return
		}
		res = rgs
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info(log.NewUrlStruct(r.URL), "fetching rss groups from mikan")
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = fmt.Errorf("fetch rss groups failed: %w", e)
		log.Error(nil, err)
	})

	c.Visit(resourcesBaseUrl + endpoint)

	return
}

func BuildSearching(s string) string {
	return resourcesBaseUrl + ResourceAPIs["search"] + s
}

// avoid nil pointer panic
func InnerTextSafety(n *html.Node) string {
	if n == nil {
		return ""
	}
	return htmlquery.InnerText(n)
}
