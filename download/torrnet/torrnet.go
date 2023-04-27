package torrnet

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	DL "github.com/NullpointerW/mikanani/download"
	"github.com/NullpointerW/mikanani/errs"
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

func DLcompl(h string) (bool, error) {
	torrs, err := DL.Qbt.TorrentList(qbt.Optional{
		"filter":  "all",
		"hashes ": h,
	})
	if err != nil {
		return false, err
	}
	if len(torrs) == 0 {
		return false, errs.Custom("%w torr hash:%s", errs.ErrTorrnetNotFound, h)
	}
	torr := torrs[0]
	return torr.Progress == 100, nil
}