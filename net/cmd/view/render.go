package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	N "github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/subject"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	"github.com/olekukonko/tablewriter"
)

type Render interface {
	RssGroup(rgs []CR.RssGroup) string
	TorrList(its []CR.Item) string
	Ls(ls []subject.Subject) string
	Status(subj *subject.Subject, torrs ...qbt.Torrent) string
}

type AsciiRender struct{}

func (_ AsciiRender) RssGroup(rgs []CR.RssGroup) string {
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

func (_ AsciiRender) TorrList(its []CR.Item) string {
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

func (_ AsciiRender) Ls(ls []subject.Subject) string {
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
	table.SetBorder(false)
	table.Render()
	return "\n" + tableString.String()
}

func (_ AsciiRender) Status(subj *subject.Subject, torrs ...qbt.Torrent) string {
	var row [][]string
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"name", "type", "resource", "file", "process", "size", "finshed", "path", "totalSize(all compl)"})
	table.SetAutoFormatHeaders(false)
	var (
		fin    string = "N"
		typ    string = subject.TV.String()
		resTyp string = subject.RSS.String()
	)
	if subj.Finished {
		fin = "Y"
	}
	if subj.Typ == subject.MOVIE {
		typ = subject.MOVIE.String()
	}
	if subj.ResourceTyp == subject.Torrent {
		resTyp = subject.Torrent.String()

	}
	var totalSize int
	for _, t := range torrs {
		fileProgress := fmt.Sprintf("%.0f", t.Progress*100) + "%"
		totalSize += t.Size
		fileSize := strconv.Itoa(t.Size/1024/1024) + "MB"
		row = append(row, []string{subj.Name, typ, resTyp, t.Name, fileProgress, fileSize, fin, subj.Path})
	}
	ttsize := strconv.Itoa(totalSize/1024/1024/1024) + "GB"
	for i, r := range row {
		row[i] = append(r, ttsize)
	}
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.AppendBulk(row)
	table.SetAutoWrapText(false)
	table.SetColWidth(60)
	table.Render()
	return "\n" + tableString.String()
}

type JsonRender struct {
	AsciiRender
}

func (_ JsonRender) RssGroup(rgs []CR.RssGroup) string {
	var rgsMap map[string][]N.TorrItem = make(map[string][]N.TorrItem)
	for _, r := range rgs {
		var nits []N.TorrItem
		for _, it := range r.Items {
			nit := N.TorrItem{}
			nit.Name = it.Name
			nit.Size = it.Size
			nit.UpdateTime = it.UpdateTime
			nits = append(nits, nit)
		}
		rgsMap[r.Name] = nits
	}
	b, _ := json.Marshal(rgsMap)
	return string(b)
}
func (_ JsonRender) TorrList(its []CR.Item) string {
	var torrls []N.TorrItem
	for _, t := range its {
		torr := N.TorrItem{}
		torr.Name = t.Name
		torr.Size = t.Size
		torr.UpdateTime = t.UpdateTime
		torrls = append(torrls, torr)
	}
	b, _ := json.Marshal(torrls)
	return string(b)
}
func (_ JsonRender) Ls(ls []subject.Subject) string {
	var sbjs []N.Subj
	for _, s := range ls {
		sbj := N.Subj{}
		sbj.Sid = s.SubjId
		sbj.Name = s.Name
		sbj.Typ = s.Typ.String()
		sbj.Compl = "N"
		if s.Terminate {
			sbj.Compl = "Y"
		}
		sbj.Status = "updating"
		if s.Finished {
			sbj.Status = "fin"
		}
		sbj.Episode = strconv.Itoa(s.Episode)
		sbjs = append(sbjs, sbj)
	}
	b, _ := json.Marshal(sbjs)
	return string(b)
}
func (jr JsonRender) Status(subj *subject.Subject, torrs ...qbt.Torrent) string {
	return jr.AsciiRender.Status(subj, torrs...)
}