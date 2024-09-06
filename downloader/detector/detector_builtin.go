package detector

import "github.com/anacrolix/torrent"

type MonitoredTorrent struct {
	Torrent *torrent.Torrent
	Rename  string
	Size    string
}
