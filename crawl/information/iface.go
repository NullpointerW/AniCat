package information

type InfoScraper interface {
	Scrape(n string) (tips map[string]string, err error)
}
