package server

import (
	"strconv"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/net/cmd/view"
	"github.com/NullpointerW/anicat/subject"
)

func route(c *cmd.Command, r view.Render) {
	switch c.Opt {
	case cmd.Add:
		sc := subject.SubjC{}
		if c.Flag.Using {
			sc.RssOption.UseRegex = c.Flag.UseRegex
			sc.RssOption.MustContain = c.Flag.MustContain
			sc.RssOption.MustNotContain = c.Flag.MustNotContain
			sc.RssOption.SubtitleGroup = c.Flag.SubtitleGroup
			sc.TorrOption.Index = c.Flag.Index
		}
		sc.N = c.N
		p := subject.NewPip(sc)
		subject.Create <- p
		c.Err = p.Error()
		if c.Err == nil {
			c.N = strconv.Itoa(p.Arg.(int))
		}
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
