package builtin

import (
	"reflect"

	util "github.com/NullpointerW/anicat/utils"
	"github.com/anacrolix/torrent"
	"golang.org/x/net/context"
)

type MonitoredTorrent struct {
	Torrent *torrent.Torrent
	Rename  string
	Size    int64
	Url     string
}
type torrentState struct {
	m       MonitoredTorrent
	gotInfo bool
}

func DetectBuiltin(recv, send chan MonitoredTorrent, ctx context.Context) {
	torrents := make(map[uintptr]torrentState)
	cases := make([]reflect.SelectCase, 0)
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx)},
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(recv)})
	for {
		c, v, _ := reflect.Select(cases)
		if c == 0 { // cancel
			return
		} else if c == 1 { // recv
			mt := v.Interface().(MonitoredTorrent)
			ts := torrentState{m: mt, gotInfo: false}
			gch := mt.Torrent.GotInfo()
			torrents[reflect.ValueOf(gch).Pointer()] = ts
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(gch)})
		} else {
			ptr := cases[c].Chan.Pointer()
			ts,ex := torrents[ptr]
			if !ex{// push ok
				goto cleancases
			}
			delete(torrents, ptr)
			if ts.gotInfo {// push
				ts.m.Size = ts.m.Torrent.Length()
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(send), Send: reflect.ValueOf(ts.m)})
			}else{// download
				ts.gotInfo=true
				ts.m.Torrent.DownloadAll()
				dch:=ts.m.Torrent.Complete.On()
				torrents[reflect.ValueOf(dch).Pointer()] = ts
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(dch)})
			}
			cleancases:
			util.SliceDelete(cases,c)
		}
	}

}
