package rss

import (
	"fmt"

	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/errs"
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

func SetAutoDLRule(rssurl, categ, dlpath, rsspath string,
	enableFilter bool, contain, exclusion string) error {
	r := qbt.AutoDLRule{
		Enabled:          true,
		AffectedFeeds:    []string{rssurl},
		SavePath:         dlpath,
		AssignedCategory: categ,
		UseRegex:         enableFilter,
		MustContain:      contain,
		MustNotContain:   exclusion,
	}
	err := DL.Qbt.SetAutoDLRule(RuleNamePrefix+rsspath, r)
	return err

}

// AddAndGetItems called only for finished items with rss type for now
func AddAndGetItems(url, path string) (*qbt.Item, error) {
	err := DL.Qbt.AddFeed(url, path)
	if err != nil {
		return nil, err
	}
	var it *qbt.Item
	ok, err := DL.DoFetch(func() (bool, error) {
		it, err = GetItems(path)
		if err != nil {
			return false, err
		}
		return len(it.Articles) > 0, nil
	}, CFG.Env.Qbt.Timeout)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: get all rss items fail,rss path=%s", errs.ErrQbtDataNotFound, path)
	}
	return it, nil
}

func GetMatchedArts(rssPath string) (arts []string, err error) {
	m, err := DL.Qbt.LsArtMatchRule(RuleNamePrefix + rssPath)
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
