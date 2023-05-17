package view

import (
	"strconv"
	"strings"

	CR "github.com/NullpointerW/mikanani/crawl/resource"
	"github.com/NullpointerW/mikanani/subject"
	"github.com/liushuochen/gotable"
	"github.com/olekukonko/tablewriter"
)

type Table interface {
	RssGroup(rgs []CR.RssGroup) string
	TorrList(its []CR.Item) string
	Ls(ls []subject.Subject) string
}

var TableRender = mergTb{}

type mergTb struct{}

func (_ mergTb) RssGroup(rgs []CR.RssGroup) string {
	var row [][]string
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"group", "name", "size", "updateTime"})
	for _, rg := range rgs {
		for _, i := range rg.Items {
			r := []string{rg.Name, i.Name, i.Size, i.UpdateTime}
			row = append(row, r)
		}
	}
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(row)
	table.SetAutoWrapText(false)
	table.SetColWidth(60)
	table.Render()
	return "\n" + tableString.String()
}

func (_ mergTb) TorrList(its []CR.Item) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"index", "name", "size", "updateTime"})
	for i, it := range its {
		table.Append([]string{strconv.Itoa(i + 1), it.Name, it.Size, it.UpdateTime})
	}
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetCenterSeparator("")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	// auto column width
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)
	table.Render()
	return "\n" + tableString.String()
}

func (_ mergTb) Ls(ls []subject.Subject) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"sid", "type", "name", "episode", "status", "compl"})
	for _, s := range ls {
		sid := strconv.Itoa(s.SubjId)
		fin := "updating"
		compl := "N"
		epi := "*"
		if s.Episode != 0 {
			epi = strconv.Itoa(s.Episode)
		}
		if s.Finished {
			fin = "fin"
		}
		if s.Terminate {
			compl = "Y"
		}
		table.Append([]string{sid, s.Typ.String(), s.Name, epi, fin, compl})
	}
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetColWidth(60)
	table.Render()
	return "\n" + tableString.String()
}

type tb struct{}

func (_ tb) RssGroup(rgs []CR.RssGroup) string {
	createItemLStb := func(its []CR.Item) string {
		tb, _ := gotable.Create("name", "size", "updateTime")
		for _, it := range its {
			tb.AddRow([]string{it.Name, it.Size, it.UpdateTime})
		}
		return "\n" + tb.String() + "\n"
	}
	ls := "\n"
	for _, rg := range rgs {
		ls += rg.Name
		ls += createItemLStb(rg.Items)
	}
	return ls
}
func (_ tb) TorrList(its []CR.Item) string {
	tb, _ := gotable.Create("index", "name", "size", "updateTime")
	for i, it := range its {
		tb.AddRow([]string{strconv.Itoa(i + 1), it.Name, it.Size, it.UpdateTime})
	}
	return "\n" + tb.String()
}

func (_ tb) Ls(ls []subject.Subject) string {
	tb, _ := gotable.Create("sid", "type", "name", "episode", "status", "compl")
	for _, s := range ls {
		sid := strconv.Itoa(s.SubjId)
		fin := "updating"
		compl := "N"
		epi := "*"
		if s.Episode != 0 {
			epi = strconv.Itoa(s.Episode)
		}
		if s.Finished {
			fin = "fin"
		}
		if s.Terminate {
			compl = "Y"
		}
		tb.AddRow([]string{sid, s.Typ.String(), s.Name, epi, fin, compl})
	}
	return "\n" + tb.String()
}
