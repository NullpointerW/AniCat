package rss

import (
	"fmt"
	RC "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/log"
	"github.com/mmcdole/gofeed"
	"regexp"
	"strings"
)

type Parser struct {
	Feed  string
	guids []string
}

const mikanbgmiIdReg = `bangumiId=(\d+)`
const mikanProjectPage = "https://mikanani.me/Home/Bangumi/%s"

func (p *Parser) GetTitleAndLink() (title string, link string, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(p.Feed)
	if err != nil {
		return
	}
	re := regexp.MustCompile(mikanbgmiIdReg)
	match := re.FindStringSubmatch(feed.Link)
	if len(match) > 1 {
		link = fmt.Sprintf(mikanProjectPage, match[1])
		link, err = RC.FetchBgmTVUrl(link)
		if err != nil {
			log.Error(log.Struct{"err", err}, "fetch bgmtv url failed")
			err = nil
			goto fetchTitle
		}
		log.Debug(nil, "link", link)
		return
	}
fetchTitle:
	title = strings.TrimPrefix(feed.Title, "Mikan Project - ")
	title = strings.TrimPrefix(title, "Mikan Project - 搜索结果:")
	return
}
