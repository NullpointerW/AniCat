package resource

import (
	"fmt"
	"log"
	"strings"

	CR "github.com/NullpointerW/mikanani/crawl"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/util"
	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
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
		a := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/ul[@class='list-inline an-ul']/li/a//@href`)
		if a != nil {
			ep, bgmurl, e := scrapeRssEndPoint(htmlquery.InnerText(a), opt)
			if e != nil {
				err = e
				return
			}
			url = resourcesBaseUrl + ep
			isrss = true
			bgmUrl = bgmurl
		} else {
			log.Println(searchstr, " rss source not found")
			if opt.Index <= 0 {
				opt.Index = 1
			}
			mglinkTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][%d]/td[1]/a[2]/@data-clipboard-text`
			mglink := htmlquery.FindOne(doc, fmt.Sprintf(mglinkTemp, opt.Index))

			if mglink != nil {
				log.Printf("\nmagnet:%s", htmlquery.InnerText(mglink))
				url = htmlquery.InnerText(mglink)
				isrss = false
			} else {
				err = errs.Custom("%w:name:%s", errs.ErrCrawlNotFound, searchstr)
			}

		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	CR.SetProxy(c)

	c.Visit(BuildSearching(CR.ConstructSearch(searchstr)))

	return
}

func scrapeRssEndPoint(endpoint string, opt Option) (rssUrl, bgmurl string, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if e != nil {
			err = e
		}
		defXpathExp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/div[@class='subgroup-text'][1]/a[@class='mikan-rss']/@href`
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
				if grpn == nil {
					grpn = htmlquery.FindOne(d, `/div[@class='dropdown']/div[@class='dropdown-toggle material-dropdown__btn']/span[1]`)
				}
				if grpn == nil {
					util.Debugf("%v not scrap group name", d)
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
		bgmXpathExp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='pull-left leftbar-container']/p[@class='bangumi-info'][last()]/a/@href`
		a := htmlquery.FindOne(doc, bgmXpathExp)
		if a == nil {
			err = errs.ErrBgmUrlNotFoundOnMikan
			return
		} else {
			bgmurl = htmlquery.InnerText(a)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.Visit(resourcesBaseUrl + endpoint)

	return
}

func ListScrape(searchstr string, t LsTyp) (res any, err error) {
	c := CR.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, e := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			err = e
			return
		}
		a := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/ul[@class='list-inline an-ul']/li/a//@href`)
		if a != nil {
			res, e = scrapeRssList(htmlquery.InnerText(a), t)
			if e != nil {
				err = e
				return
			}
		} else {
			log.Println(searchstr, "resource ls: torr typ")
			switch t {
			case Ls:
				fnTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][%d]/td[1]/a[@class='magnet-link-wrap']`
				szTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][%d]/td[2]`
				uptTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][%d]/td[3]`
				// torrTemp := htmlquery.FindOne(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][%d]/td[4]/a/@href`)
				nodes := htmlquery.Find(doc, `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row']`)
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
					err = errs.Custom("%w: %s %s", errs.ErrCrawlNotFound, t.String(), searchstr)
				} else {
					res, err = items, nil
				}
			case LSGroup:
				err = errs.Custom("%w:search item `%s` is torrent type ", errs.ErrLsGroupUnavailableOnTorr, searchstr)
			default:
				err = errs.ErrUnknownResCrawlLsType
			}

		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	CR.SetProxy(c)

	c.Visit(BuildSearching(CR.ConstructSearch(searchstr)))

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
		trsTemp := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/table[@class='table table-striped tbl-border fadeIn'][%d]/tbody/tr`
		itnExp := `/td[1]/a[@class='magnet-link-wrap']`
		szExp := `/td[2]`
		uptExp := `/td[3]`
		var rgs []RssGroup
		tg := `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/div[@class='subgroup-text']`
		ds := htmlquery.Find(doc, tg)
		for i, d := range ds {
			rg := RssGroup{}
			grpn := htmlquery.FindOne(d, `/a[1]`)
			if grpn == nil {
				grpn = htmlquery.FindOne(d, `/div[@class='dropdown']/div[@class='dropdown-toggle material-dropdown__btn']/span[1]`)
			}
			if grpn == nil {
				util.Debugf("%v not scrap group name", d)
				continue
			}
			rg.Name = htmlquery.InnerText(grpn)
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
			err = errs.ErrCrawlNotFound
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
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		err = e
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
