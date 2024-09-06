package builtin

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

var DefaultDownLoader *Downloader

func InitDownloader() {
	DefaultDownLoader = NewDownloader("./.db", false, NewHttpSeeker())
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
	ts.Storage = storage.NewFileOpts(fop)
	t, _, err = d.client.AddTorrentSpec(ts)
	return
}
