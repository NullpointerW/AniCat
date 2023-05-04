package rss

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	DL "github.com/NullpointerW/mikanani/download"
)



func Download(adlr qbt.AutoDLRule, path string) (err error) {
	err = DL.Qbt.AddFeed(adlr.AffectedFeeds[0], path)
	if err != nil {
		return err
	}
	err = DL.Qbt.SetAutoDLRule("ADL-"+path, adlr)
	return err
}
