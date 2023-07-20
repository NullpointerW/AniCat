package server

import (
	"strconv"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/net/cmd"
	"github.com/NullpointerW/anicat/net/cmd/view"
	"github.com/NullpointerW/anicat/subject"
	util "github.com/NullpointerW/anicat/utils"
)

func route(c *cmd.Command) {
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
	case cmd.Del:
		i, err := strconv.Atoi(c.N)
		if err != nil {
			c.Err = err
			return
		}
		p := subject.NewPip(i)
		subject.Delete <- p
		c.Err = p.Error()
	case cmd.Ls:
		util.Debugln("at cmd.Ls")
		ls := subject.Manager.List()
		c.N = view.TableRender.Ls(ls)
		// util.Debugln("render table ::", c.N)
	case cmd.LsItems:
		l, err := CR.ListScrape(c.N, CR.Ls)
		if err != nil {
			c.Err = err
			return
		}
		ls := ""
		rgs, RssGroupSlice := l.([]CR.RssGroup)
		its, ItemSlice := l.([]CR.Item)
		if RssGroupSlice {
			ls = view.TableRender.RssGroup(rgs)
		} else if ItemSlice {
			ls = view.TableRender.TorrList(its)
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
			c.N = view.TableRender.Status(subj, h)
		} else {
			hs, err := torrent.GetViaCateg(subj.QbtCateg())
			if err != nil {
				c.Err = err
				return
			}
			c.N = view.TableRender.Status(subj, hs...)
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
