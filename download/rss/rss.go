package rss

import (
	"time"

	DL "github.com/NullpointerW/anicat/download"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
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

func SetAutoDLRule(rssurl, categ, dlpath, rsspath string) error {
	r := qbt.AutoDLRule{
		Enabled:          true,
		AffectedFeeds:    []string{rssurl},
		SavePath:         dlpath,
		AssignedCategory: categ,
	}
	err := DL.Qbt.SetAutoDLRule(RuleNamePrefix+rsspath, r)
	return err

}

func AddAndGetItems(url, path string) (*qbt.Item, error) {
	err := DL.Qbt.AddFeed(url, path)
	if err != nil {
		return nil, err
	}
	time.Sleep(1500 * time.Millisecond)
	it, err := GetItems(path)
	return it, err
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

func GetItems(rssPath string) (*qbt.Item, error) {
	its, err := DL.Qbt.GetAllItems(true)
	if err != nil {
		return nil, err
	}
	if it, e := its[rssPath]; !e {
		return nil, nil
	} else {
		return &it, nil
	}

}

func RmRss(rssPath string) error {
	err := DL.Qbt.RemoveItem(rssPath)
	if err != nil {
		return err
	}
	adlrn := RuleNamePrefix + rssPath
	return DL.Qbt.RmAutoDLRule(adlrn)
}
