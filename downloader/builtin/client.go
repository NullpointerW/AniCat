package builtin

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

var DefaultDownLoader *Downloader

func init() {
	if CFG.Env.BuiltinDownloader {
		log.Info(log.Struct{"github", "https://github.com/anacrolix/torrent"}, "builtin-enabled, using anacrolix/torrent")
		InitDownloader()
	}
}
func InitDownloader() {
	cfg := &DownloaderConfig{
		BaseDir:    "./.db",
		FakePeerID: true,
		NopUpload:  true,
		Seeker:     NewHttpSeeker(),
	}
	DefaultDownLoader = NewDownloader(cfg)
}

type Downloader struct {
	client *torrent.Client
	TorrentSeeker
	extraTrackers [][]string
}
type FileName interface {
	Name() storage.FilePathMaker
}

type FileOption interface {
	FileName
	Dir() storage.TorrentDirFilePathMaker
}
type DownloaderConfig struct {
	BaseDir    string
	Seeker     TorrentSeeker
	NopUpload  bool
	FakePeerID bool
}

func NewDownloader(c *DownloaderConfig) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = !c.NopUpload
	cfg.NoUpload = c.NopUpload
	fop := storage.NewFileClientOpts{}
	fop.ClientBaseDir = c.BaseDir
	cfg.DefaultStorage = storage.NewFileOpts(fop)
	if c.Seeker == nil {
		c.Seeker = NewHttpSeeker()
	}
	if c.FakePeerID {
		f := "-qB419E-" // qBittorrent
		var b [20]byte
		n := copy(b[:], ([]byte)(f))
		_, err := rand.Read(b[n:])
		if err != nil {
			panic("builtin-downloader: error generating peer id")
		}
		cfg.PeerID = (string)(b[:])
		cfg.HTTPUserAgent = "qBittorrent/v4.1.9.14"
		mainPath := "github.com/NullpointerW/anicat"
		mainVersion := "unknown"
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			mainPath = buildInfo.Main.Path
			mainVersion = buildInfo.Main.Version
		}
		exhskVer := fmt.Sprintf(
			"%v %v (%v %v)",
			mainPath,
			mainVersion,
			"qBittorrent",
			"v4.1.9.14",
		)
		cfg.ExtendedHandshakeClientVersion = exhskVer
	}
	client, err := torrent.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	peerIDStr := func(p torrent.PeerID) string {
		b := ([20]byte)(p)
		// fmt.Println((string)(b[:]))
		return string(b[:])
	}

	log.Info(log.Struct{"version", cfg.ExtendedHandshakeClientVersion, "userAgent", cfg.HTTPUserAgent, "peerID", peerIDStr(client.PeerID()), "upnpID", cfg.UpnpID}, "torrent-client initialized")
	return &Downloader{client, c.Seeker, extraTrackers()}

}

func (d *Downloader) Download(s string, fOp FileOption, seeker TorrentSeeker) (t *torrent.Torrent, err error) {
	var ts *torrent.TorrentSpec
	if seeker == nil {
		ts, err = d.Seek(s)
	} else {
		ts, err = seeker.Seek(s)
	}
	if err != nil {
		return nil, err
	}

	fop := storage.NewFileClientOpts{}
	fop.TorrentDirMaker = fOp.Dir()
	fop.FilePathMaker = fOp.Name()
	fop.PieceCompletion = storage.NewMapPieceCompletion()
	ts.Storage = storage.NewFileOpts(fop)
	t, _, err = d.client.AddTorrentSpec(ts)
	t.AddTrackers(d.extraTrackers)
	return
}
func extraTrackers() (trackers [][]string) {
	provider := "https://cdn.jsdelivr.net/gh/DeSireFire/animeTrackerList/AT_all.txt"
	h := http.Client{Timeout: time.Second * 5}
	get, err := h.Get(provider)
	if err != nil {
		log.Error(log.Struct{"err", err}, "builtin-downloader: set tracker failed")
		return nil
	}
	defer get.Body.Close()
	r := bufio.NewScanner(get.Body)
	var trackerURLs []string
	for r.Scan() {
		trackerURLs = append(trackerURLs, r.Text())
	}
	return append(trackers, trackerURLs)
}
