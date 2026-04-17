package view

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	CR "github.com/NullpointerW/anicat/crawl/resource"
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/log"
	"github.com/NullpointerW/anicat/net"
	"github.com/NullpointerW/anicat/subject"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

type TorrItem struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime string `json:"uptime"`
}

type Subj struct {
	Sid     int    `json:"sid"`
	Typ     string `json:"type"`
	Name    string `json:"name"`
	Episode string `json:"epi"`
	Status  string `json:"status"`
	Compl   string `json:"compl"`
}

type Stat struct {
	File     string `json:"file"`
	Size     string `json:"size"`
	Progress string `json:"progress"`
}

type StatPayload struct {
	Path      string `json:"path"`
	TotalSize string `json:"total_size"`
	Torrents  []Stat `json:"torrents"`
}

type JsonRender struct {
	Conn *net.Conn
}

func (JsonRender) RssGroup(rgs []CR.RssGroup) string {
	rgsMap := make(map[string][]TorrItem)
	for _, r := range rgs {
		var nits []TorrItem
		for _, it := range r.Items {
			nits = append(nits, TorrItem{Name: it.Name, Size: it.Size, UpdateTime: it.UpdateTime})
		}
		rgsMap[r.Name] = nits
	}
	b, _ := json.Marshal(rgsMap)
	return string(b)
}

func (JsonRender) TorrList(its []CR.Item) string {
	var torrls []TorrItem
	for _, t := range its {
		torrls = append(torrls, TorrItem{Name: t.Name, Size: t.Size, UpdateTime: t.UpdateTime})
	}
	b, _ := json.Marshal(torrls)
	return string(b)
}

func (JsonRender) Ls(ls []subject.Subject) string {
	var sbjs []Subj
	for _, s := range ls {
		sbj := Subj{
			Sid:     s.SubjId,
			Name:    s.Name,
			Typ:     s.Typ.String(),
			Compl:   "N",
			Status:  "updating",
			Episode: strconv.Itoa(s.Episode),
		}
		if s.Terminate {
			sbj.Compl = "Y"
		}
		if s.Finished {
			sbj.Status = "fin"
		}
		sbjs = append(sbjs, sbj)
	}
	b, _ := json.Marshal(sbjs)
	return string(b)
}

func (JsonRender) Status(subj *subject.Subject, torrs ...qbt.Torrent) string {
	var totalSize int
	var stats []Stat
	for _, t := range torrs {
		totalSize += t.Size
		stats = append(stats, Stat{
			File:     filepath.Base(t.ContentPath),
			Size:     strconv.Itoa(t.Size/1024/1024) + "MB",
			Progress: fmt.Sprintf("%.0f", t.Progress*100) + "%",
		})
	}
	payload := StatPayload{
		Path:      subj.Path,
		TotalSize: strconv.Itoa(totalSize/1024/1024/1024) + "GB",
		Torrents:  stats,
	}
	b, _ := json.Marshal(payload)
	return string(b)
}

func (r JsonRender) StatusBuiltin(subj *subject.Subject) {
	r.Conn.Hajacked = true
	go HandleStatus(subj, r.Conn)
}

func HandleStatus(s *subject.Subject, c *net.Conn) {
	c.Write("keep-alive")
	var list builtin.TorrentProgressList
	for {
		if s.Terminate {
			<-s.Exited
			list.Put(s.FinishedTorrentNameList.List())
			r, _ := json.Marshal(builtin.TorrentProgressListSend{List: list.Get(), Fin: true})
			if err := c.Write(string(r)); err != nil {
				log.Error(log.Struct{"err", err}, "conn write err")
			}
			c.TcpConn.Close()
			return
		}
		list.Put(s.FinishedTorrentNameList.List())
		list.Put(s.TorrentMonitor.GetProgressList())
		lse := builtin.TorrentProgressListSend{List: list.Get(), Fin: list.Fin()}
		r, _ := json.Marshal(lse)
		if err := c.Write(string(r)); err != nil {
			log.Error(log.Struct{"err", err}, "conn write err")
			return
		}
		if lse.Fin {
			c.TcpConn.Close()
			return
		}
		time.Sleep(15 * time.Second)
	}
}
