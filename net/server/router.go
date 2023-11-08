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

func route(c cmd.Cmd, r view.Render) (resp string, err error) {
	switch c.Cmd {
	case cmd.Add:
		return addSubjProcess(c, false)
	case cmd.AddFeed:
		return addSubjProcess(c, true)
	case cmd.Remove:
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
	case cmd.Ls:
		ls := subject.Manager.List()
		resp = r.Ls(ls)
		return
	case cmd.LsItems:
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
	case cmd.Status:
		i, er := strconv.Atoi(c.Arg)
		if er != nil {
			err = er
			return
		}
		subj := subject.Manager.Get(i)
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

	case cmd.Stop:
		for _, s := range subject.Manager.List() {
			if !s.Terminate {
				s.Exit()
			}
		}
		resp = "exited."
		return
	default:
		return
	}
}

func addSubjProcess(c cmd.Cmd, isFeed bool) (resp string, err error) {
	flag := new(cmd.AddFlag)
	err = json.Unmarshal(c.Raw, &flag)
	if err != nil {
		return "", err
	}
	sc := transferSubjC(c.Arg, bytes.Equal(c.Raw, []byte("null")), *flag, isFeed)
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
