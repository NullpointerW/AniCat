package builtin

import (
	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

var DefaultDownLoader *Downloader

func init() {
	if CFG.Env.BuiltinDownloader {
		InitDownloader()
		log.Info(log.Struct{"github", "https://github.com/anacrolix/torrent"}, "builtin downloader enabled, using anacrolix/torrent")
	}
}
func InitDownloader() {
	DefaultDownLoader = NewDownloader("./.db", true, NewHttpSeeker())
}

type Downloader struct {
	client *torrent.Client
	TorrentSeeker
}
type FileName interface {
	Name() storage.FilePathMaker
}

type FileOption interface {
	FileName
	Dir() storage.TorrentDirFilePathMaker
}

func NewDownloader(basedir string, nopUpload bool, seeker TorrentSeeker) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = !nopUpload
	fop := storage.NewFileClientOpts{}
	fop.ClientBaseDir = basedir
	cfg.DefaultStorage = storage.NewFileOpts(fop)
	client, err := torrent.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	log.Info(log.Struct{"version",cfg.ExtendedHandshakeClientVersion,"userAgent",cfg.HTTPUserAgent,"peerID",cfg.PeerID,"upnpID",cfg.UpnpID},"torrent-client initialized")
	return &Downloader{client, seeker}

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
	return
}
