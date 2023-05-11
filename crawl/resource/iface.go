package resource

type ResScraper interface {
	Scrape(searchstr string) (url, bgmUrl string, isrss bool, err error)
}
