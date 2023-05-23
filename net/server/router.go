package server

import (
	// "fmt"
	"strconv"

	// "strings"

	CR "github.com/NullpointerW/mikanani/crawl/resource"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/net/cmd"
	"github.com/NullpointerW/mikanani/net/cmd/view"
	"github.com/NullpointerW/mikanani/subject"
	"github.com/liushuochen/gotable"
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
		ls := subject.Manager.List()
		c.N = view.TableRender.Ls(ls)
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

	case cmd.SavePath:
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
		c.N = subj.Path
	}
}

func createItemLStb(its []CR.Item) string {
	tb, _ := gotable.Create("name", "size", "updateTime")
	for _, it := range its {
		tb.AddRow([]string{it.Name, it.Size, it.UpdateTime})
	}
	return "\n" + tb.String() + "\n"
}
