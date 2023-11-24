package downloader

type (
	Torrent interface {
		Add(magnet, path string)
	}
	Rss interface {
		GetTorrents() ([]string, error)
	}
	Downloader interface {
		Torrent
		Rss
	}
)
