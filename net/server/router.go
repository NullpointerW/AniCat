package server

import (
	"bytes"
	"encoding/json"
	"strconv"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/net/cmd/view"
	"github.com/NullpointerW/anicat/subject"
)

var CmdSelector *cmd.Selector

func init() {
	add := cmd.NewCommandCase(cmd.Add, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		return addSubjProcess(c, false)
	})
	addFeed := cmd.NewCommandCase(cmd.AddFeed, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		return addSubjProcess(c, true)
	})
	remove := cmd.NewCommandCase(cmd.Remove, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		i, er := strconv.Atoi(c.Arg)
		if er != nil && c.Arg != "*" {
			err = er
			return
		}
		var pip *subject.Pip
		if c.Arg == "*" {
			pip = subject.NewPip("*")
		} else {
			pip = subject.NewPip(i)
		}
		subject.Delete <- pip
		return "ok", pip.Error()
	})
	list := cmd.NewCommandCase(cmd.Ls, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		ls := subject.Mgr.List()
		resp = r.Ls(ls)
		return
	})
	listItem := cmd.NewCommandCase(cmd.LsItems, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		flag := new(cmd.LsiFlag)
		err = json.Unmarshal(c.Raw, &flag)
		if err != nil {
			return "", err
		}
		search := false
		if !bytes.Equal(c.Raw, []byte("null")) {
			search = flag.SearchList
		}
		l, er := CR.ListScrape(c.Arg, CR.Ls, search)
		if er != nil {
			err = er
			return
		}
		ls := ""
		rgs, RssGroupSlice := l.([]CR.RssGroup)
		its, ItemSlice := l.([]CR.Item)
		if RssGroupSlice {
			ls = r.RssGroup(rgs)
		} else if ItemSlice {
			ls = r.TorrList(its)
		} else {
			err = errs.ErrUndefinedCrawlListType
		}
		resp = ls
		return
	})
	status := cmd.NewCommandCase(cmd.Status, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		i, er := strconv.Atoi(c.Arg)
		if er != nil {
			err = er
			return
		}
		subj := subject.Mgr.Get(i)
		if subj == nil {
			err = errs.ErrSubjectNotFound
			return
		}

		if subj.ResourceTyp == subject.Torrent {
			h, er := torrent.Get(subj.TorrentHash)
			if er != nil {
				err = er
				return
			}
			resp = r.Status(subj, h)
		} else {
			hs, er := torrent.GetViaCateg(subj.QbtCateg())
			if er != nil {
				err = er
				return
			}
			resp = r.Status(subj, hs...)
		}
		return
	})
	stop := cmd.NewCommandCase(cmd.Stop, func(c cmd.Cmd, r view.Render) (resp string, err error) {
		for _, s := range subject.Mgr.List() {
			if !s.Terminate {
				s.Exit()
			}
		}
		subject.Mgr.Exit()
		resp = "exited."
		return
	})
	CmdSelector = cmd.NewSelector(add, addFeed, remove, list, listItem, status, stop)
}

func route(c cmd.Cmd, r view.Render) (resp string, err error) {
	return CmdSelector.Select(c, r)
}

func addSubjProcess(c cmd.Cmd, isFeed bool) (resp string, err error) {
	flag := new(cmd.AddFlag)
	err = json.Unmarshal(c.Raw, &flag)
	if err != nil {
		return "", err
	}
	sc := transferSubjC(c.Arg, !bytes.Equal(c.Raw, []byte("null")), *flag, isFeed)
	resp, err = createSubject(sc)
	return
}
func transferSubjC(arg string, usingFlag bool, src cmd.AddFlag, feed bool) (dst subject.SubjC) {
	if usingFlag {
		dst.RssOption.UseRegex = src.UseRegexp
		dst.RssOption.MustContain = src.MustContain
		dst.RssOption.MustNotContain = src.MustNotContain
		dst.RssOption.SubtitleGroup = src.Group
		dst.TorrOption.Index = src.Index
		dst.RssOption.Name = src.FeedInfoName
	}
	dst.N = arg
	dst.CreateTyp = subject.CreateViaStr
	if feed {
		dst.CreateTyp = subject.CreateViaFeed
	}
	return dst
}

func createSubject(sc subject.SubjC) (resp string, err error) {
	p := subject.NewPip(sc)
	subject.Create <- p
	err = p.Error()
	if err == nil {
		resp = strconv.Itoa(p.Arg.(int))
	}
	return
}
