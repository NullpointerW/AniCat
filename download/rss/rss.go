package rss

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	DL "github.com/NullpointerW/mikanani/download"
)

const RuleNamePrefix = "ADL-"

func Download(adlr qbt.AutoDLRule, path string) (err error) {
	err = DL.Qbt.AddFeed(adlr.AffectedFeeds[0], path)
	if err != nil {
		return err
	}
	err = DL.Qbt.SetAutoDLRule(RuleNamePrefix+path, adlr)
	return err
}

func GetMatchedArts(rssPath string) (arts []string, err error) {
	m, err := DL.Qbt.LsArtMatchRlue(RuleNamePrefix + rssPath)
	if err != nil {
		return nil, err
	}
	for _, v := range m {
		arts = append(arts, v...)
	}
	return arts, nil
}

func RmRss(rssPath string) error {
	err := DL.Qbt.RemoveItem(rssPath)
	if err != nil {
		return err
	}
	adlrn := RuleNamePrefix + rssPath
	return DL.Qbt.RmAutoDLRule(adlrn)
}
