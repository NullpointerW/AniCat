package downloader

type (
	Torrent interface {
		Add(magnet, path string)
		Process()
	}
	Rss interface {
		GetTorrents() ([]string, error)
	}
	Downloader interface {
		Torrent
		Rss
	}
)
