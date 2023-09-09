package server

import (
	"strconv"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/net/cmd/view"
	"github.com/NullpointerW/anicat/subject"
)

func route(c *cmd.Command, r view.Render) {
	switch c.Opt {
	case cmd.Add:
		sc := transferSubjC(c, false)
		createSubject(c, sc)
	case cmd.AddFeed:
		sc := transferSubjC(c, true)
		createSubject(c, sc)
	case cmd.Del:
		i, err := strconv.Atoi(c.N)
		if err != nil && c.N != "*" {
			c.Err = err
			return
		}
		var pip *subject.Pip
		if c.N == "*" {
			pip = subject.NewPip("*")
		} else {
			pip = subject.NewPip(i)
		}
		subject.Delete <- pip
		c.Err = pip.Error()
	case cmd.Ls:
		ls := subject.Manager.List()
		c.N = r.Ls(ls)
	case cmd.LsItems, cmd.LsItems_searchlist:
		isLsi_l := c.Opt == cmd.LsItems_searchlist
		l, err := CR.ListScrape(c.N, CR.Ls, isLsi_l)
		if err != nil {
			c.Err = err
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
			c.Err = errs.ErrUndefinedCrawlListType
		}
		c.N = ls
	case cmd.LsGroup:
		c.Err = errs.WarnReservedCommand_lsg

	case cmd.Status:
		i, err := strconv.Atoi(c.N)
		if err != nil {
			c.Err = err
			return
		}
		subj := subject.Manager.Get(i)

		if subj == nil {
			c.Err = errs.ErrSubjectNotFound
			return
		}

		if subj.ResourceTyp == subject.Torrent {
			h, err := torrent.Get(subj.TorrentHash)
			if err != nil {
				c.Err = err
				return
			}
			c.N = r.Status(subj, h)
		} else {
			hs, err := torrent.GetViaCateg(subj.QbtCateg())
			if err != nil {
				c.Err = err
				return
			}
			c.N = r.Status(subj, hs...)
		}

	case cmd.Stop:
		for _, s := range subject.Manager.List() {
			if !s.Terminate {
				s.Exit()
			}
		}
		c.N = "exited."
	}
}
func transferSubjC(src *cmd.Command, feed bool) (dst subject.SubjC) {
	if src.Flag.Using {
		dst.RssOption.UseRegex = src.Flag.UseRegex
		dst.RssOption.MustContain = src.Flag.MustContain
		dst.RssOption.MustNotContain = src.Flag.MustNotContain
		dst.RssOption.SubtitleGroup = src.Flag.SubtitleGroup
		dst.TorrOption.Index = src.Flag.Index
		dst.RssOption.Name = src.Flag.Name
	}
	dst.N = src.N
	dst.CreateTyp = subject.CreateViaStr
	if feed {
		dst.CreateTyp = subject.CreateViaFeed
	}
	return dst
}

func createSubject(c *cmd.Command, sc subject.SubjC) {
	p := subject.NewPip(sc)
	subject.Create <- p
	c.Err = p.Error()
	if c.Err == nil {
		c.N = strconv.Itoa(p.Arg.(int))
	}
}
