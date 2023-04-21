package crawl

import (
	CFG "github.com/NullpointerW/mikanani/conf"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func setProxy(c *colly.Collector) {
	if CFG.Proxy != nil {
		if p, err := proxy.RoundRobinProxySwitcher(
			CFG.Proxy...,
		); err == nil {
			c.SetProxyFunc(p)
		}
	}
}
