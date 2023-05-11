package torrent

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	DL "github.com/NullpointerW/mikanani/download"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/util"
)

func Add(url, path, tag string) (string, error) {
	err := DL.Qbt.AddNewTorrentViaUrl(url, path, tag)
	if err != nil {
		return "", err
	}

	torrs, err := DL.Qbt.TorrentList(qbt.Optional{
		"filter": "all",
		"tag":    tag,
	})
	if err != nil {
		return "", err
	}

	return torrs[0].Hash, nil
}

func Get(h string) (torr qbt.Torrent, err error) {
	torrs, err := DL.Qbt.TorrentList(qbt.Optional{
		"filter": "all",
		"hashes": h,
	})
	if err != nil {
		return torr, err
	}
	if len(torrs) == 0 {
		return torr, errs.Custom("%w:torr hash:%s", errs.ErrTorrnetNotFound, h)
	}
	torr = torrs[0]
	return
}

func GetViaPath(path string) (hits []qbt.Torrent, err error) {
	path=util.FileSeparatorConv(path)
	torrs, err := DL.Qbt.TorrentList(qbt.Optional{
		"filter": "all",
	})
	if err != nil {
		return hits, err
	}
	if len(torrs) == 0 {
		return hits, errs.Custom("%w:there is no any torrents on qbt", errs.ErrTorrnetNotFound)
	}

	for _, t := range torrs {
		p := util.FileSeparatorConv(t.SavePath)
		util.Debugln("torr_save_path:",p)
		if p == path {
			hits = append(hits, t)
		}
	}
	if len(hits) == 0 {
		return hits, errs.Custom("%w:no torrents found for \"%s\" wtih save path", errs.ErrTorrnetNotFound, path)
	}
	return
}

func DLcompl(h string) (bool, error) {
	torr, err := Get(h)
	if err != nil {
		return false, err
	}
	return torr.Progress == 1, nil
}

func DelTorrs(p string) error {
	// util.Debugln("abs_path:", p)
	torrs, err := GetViaPath(p)
	if err != nil {
		return err
	}
	var hs []string
	for _, torr := range torrs {
		hs = append(hs, torr.Hash)
	}
	return DL.Qbt.DelTorrentsFs(hs...)
}

func DelTag(t string) error {
	return DL.Qbt.DelTags(t)
}
