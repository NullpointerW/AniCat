package builtin

import (
	"reflect"

	"github.com/anacrolix/torrent"
	"golang.org/x/net/context"
)

type MonitoredTorrent struct {
	Torrent *torrent.Torrent
	Rename  string
	Size    string
	Url     string
}
type torrentState struct {
	t       MonitoredTorrent
	gotInfo bool
}

func DetectBuiltin(recv, send chan MonitoredTorrent, ctx context.Context) {
	torrents := make(map[uintptr]torrentState)
	cases := make([]reflect.SelectCase, 0)
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault},
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx)},
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(recv)})
	for {
		c, _, _ := reflect.Select(cases)
		if c == 0 { // default

		}
	}

}
