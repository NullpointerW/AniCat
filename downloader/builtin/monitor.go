package builtin

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	util "github.com/NullpointerW/anicat/utils"
	"github.com/anacrolix/torrent"
	"golang.org/x/net/context"
)

type TorrentProgressList struct {
	list  map[string]int
	clean []TorrentProgress
}

func (l *TorrentProgressList) Put(dirty []TorrentProgress) {
	if l.list == nil {
		l.list = make(map[string]int)
	}
	for _, t := range dirty {
		if idx, ex := l.list[t.Name]; !ex {
			l.list[t.Name] = len(l.clean)
			l.clean = append(l.clean, t)
		} else {
			l.clean[idx] = t
		}
	}
}
func (l *TorrentProgressList) Get() []TorrentProgress {
	if l.clean == nil {
		l.clean = make([]TorrentProgress, 0)
	}
	return l.clean
}
func (l *TorrentProgressList) Fin() bool {
	for _, t := range l.clean {
		if t.Percentage != 100 {
			return false
		}
	}
	return true
}

type TorrentProgressListSend struct {
	List []TorrentProgress `json:"list"`
	Fin  bool              `json:"fin"`
}
type TorrentProgress struct {
	Percentage int    `json:"percentage"`
	Name       string `json:"name"`
}

type TorrentProgressMonitor struct {
	mu              sync.Mutex
	activeTorrents  map[TorrentInfo]struct{}
	cache, lasttime atomic.Value
	ttl             time.Duration
}

func NewTorrentProgressMonitor(ttl time.Duration) *TorrentProgressMonitor {
	tm := &TorrentProgressMonitor{ttl: ttl, activeTorrents: make(map[TorrentInfo]struct{})}
	return tm
}

func (tm *TorrentProgressMonitor) AddTorrent(t TorrentInfo) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.activeTorrents[t] = struct{}{}
}

func (tm *TorrentProgressMonitor) GetProgressList() []TorrentProgress {
	if c, ok := tm.checkAndGet(); ok {
		return c
	}
	return tm.getSlow()
}
func (tm *TorrentProgressMonitor) getSlow() []TorrentProgress {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if c, ok := tm.checkAndGet(); ok {
		return c
	}
	c := make([]TorrentProgress, 0, len(tm.activeTorrents))
	for t := range tm.activeTorrents {
		tt := t.Torrent
		p := int(float64(tt.BytesCompleted()) / float64(tt.Length()) * 100)
		if p == 100 {
			delete(tm.activeTorrents, t)
		}
		tg := TorrentProgress{Percentage: p, Name: t.Rename}
		c = append(c, tg)
	}
	tm.cache.Store(c)
	tm.lasttime.Store(time.Now())
	return c
}

func (tm *TorrentProgressMonitor) checkAndGet() ([]TorrentProgress, bool) {
	if l, lt, now := tm.cache.Load(), tm.lasttime.Load(), time.Now(); lt != nil && lt.(time.Time).Add(tm.ttl).Sub(now) >= 0 && l != nil {
		return l.([]TorrentProgress), true
	}
	return nil, false
}

type TorrentInfo struct {
	Torrent *torrent.Torrent
	Rename  string
}
type MonitoredTorrent struct {
	TorrentInfo
	Size int64
	Url  string
}
type torrentState struct {
	MonitoredTorrent
	gotInfo bool
}

// MonitorBuiltin monitors torrents and handles their state transitions.
// It listens for incoming MonitoredTorrent objects on the recv channel,
// processes them, and sends updated MonitoredTorrent objects on the send channel.
// The function also listens for context cancellation to gracefully exit.
//
// Parameters:
// - recv: a channel from which MonitoredTorrent objects are received.
// - send: a channel to which updated MonitoredTorrent objects are sent.
// - ctx: a context to handle cancellation and timeout.
//
// The function uses reflection to dynamically select from multiple channels
// and manage the state of each torrent. It tracks torrent states in a map
// and updates their status based on events such as receiving torrent info
// or completing the download.
func MonitorBuiltin(recv, send chan MonitoredTorrent, ctx context.Context, monitor *TorrentProgressMonitor) {
	torrents := make(map[uintptr]torrentState)
	cases := make([]reflect.SelectCase, 0)
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(recv)})
	for {
		c, v, _ := reflect.Select(cases)
		if c == 0 { // cancel
			close(recv)
			close(send)
			for _, t := range torrents {
				t.Torrent.Drop()
			}
			return
		} else if c == 1 { // recv
			mt := v.Interface().(MonitoredTorrent)
			ts := torrentState{MonitoredTorrent: mt, gotInfo: false}
			gch := mt.Torrent.GotInfo()
			torrents[reflect.ValueOf(gch).Pointer()] = ts
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(gch)})
			fmt.Printf("bultin-detector: recv torrent download-event,start get info:%+v \n", mt)
		} else {
			ptr := cases[c].Chan.Pointer()
			ts, ex := torrents[ptr]
			if !ex { // push ok
				goto cleancases
			}
			delete(torrents, ptr)
			if ts.gotInfo { // push
				ts.Size = ts.Torrent.Length()
				fmt.Printf("bultin-detector: download complete :%+v \n", ts.MonitoredTorrent)
				ts.Torrent.Drop()
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(send), Send: reflect.ValueOf(ts.MonitoredTorrent)})
			} else { // download
				ts.gotInfo = true
				ts.Torrent.DownloadAll()
				dch := ts.Torrent.Complete.On()
				if ts.Rename == "" {
					ts.Rename = ts.Torrent.Name()
				}
				torrents[reflect.ValueOf(dch).Pointer()] = ts
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(dch)})
				monitor.AddTorrent(ts.TorrentInfo)
				fmt.Printf("bultin-detector: got torrent info ok,downloading :%+v \n", ts.MonitoredTorrent)
			}
		cleancases:
			cases = util.SliceDelete(cases, c)
		}
	}

}
