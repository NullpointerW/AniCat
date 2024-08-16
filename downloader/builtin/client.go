package builtin

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

var DefaultDownLoader *Downloader

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
	cfg.NoUpload = nopUpload
	fop := storage.NewFileClientOpts{}
	fop.ClientBaseDir = basedir
	cfg.DefaultStorage = storage.NewFileOpts(fop)
	client, err := torrent.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &Downloader{client, seeker}
}

func (d *Downloader) Download(s string, fOp FileOption) (*torrent.Torrent, error) {
	reader, err := d.Seek(s)
	if err != nil {
		return nil, err
	}
	mf, err := metainfo.Load(reader)
	if err != nil {
		return nil, err
	}
	ts := torrent.TorrentSpecFromMetaInfo(mf)
	fop := storage.NewFileClientOpts{}
	fop.TorrentDirMaker = fOp.Dir()
	fop.FilePathMaker = fOp.Name()
	ts.Storage = storage.NewFileOpts(fop)
	t, _, err := d.client.AddTorrentSpec(ts)
	return t, err

}
