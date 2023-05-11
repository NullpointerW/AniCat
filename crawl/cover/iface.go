package cover

type CoverScraper interface {
	Scrape(filePath, CoverName string) error
}

type CoverScraperFunc func(string, string) error

func (f CoverScraperFunc) Scrape(filePath, CoverName string) error {
	return f(filePath, CoverName)
}

var DOUBANCoverScraper = CoverScraper(CoverScraperFunc(TouchCoverImg))