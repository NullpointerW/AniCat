package information

import sel "github.com/NullpointerW/anicat/crawl/selector"

const (
	infoBaseUrl = `https://bgm.tv`

	// fallback XPaths used when selectors.yaml is not loaded
	infoPageXpathExpDefault = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@id='main'][2]/div[@class='columns clearit']/div[@id='columnSearchB']/ul[@id='browserItemList']/li[1]/div[@class='inner']/h3/a[@class='l']/@href`
	infoXpathExpDefault     = `//ul[@id='infobox']/li`
	originNameXpathDefault  = `//h1[@class='nameSingle']/a`
)

// keys for info map Scrape from bgm.tv
const (
	SubjId            = "sid"
	SubjName          = "中文名"
	SubjOriginName    = "originName"
	SubjEpisode       = "话数"
	SubjStartTime     = "放送开始"
	SubjMoveStartTime = "上映年度"
	SubjectEndTime    = "播放结束"
	Alias             = "别名"
)

var InfoAPIs = map[string]string{
	"search":  "/subject_search/%s?cat=2",
	"subject": "/subject/%d",
}

var TMDBAPIs = map[string]string{
	"search": "/search/%s?query=%s",
}

const (
	TMDB_HOST      = "https://www.themoviedb.org"
	TMDB_TYP_TV    = "tv"
	TMDB_TYP_MOVIE = "movie"
)

func bgmXpath(key, fallback string) string {
	if v := sel.Get("bgmtv", key); v != "" {
		return v
	}
	return fallback
}

func tmdbXpath(key, fallback string) string {
	if v := sel.Get("tmdb", key); v != "" {
		return v
	}
	return fallback
}
