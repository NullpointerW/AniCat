package test

import (
	"fmt"
	DL "github.com/NullpointerW/anicat/downloader"
	TORR "github.com/NullpointerW/anicat/downloader/torrent"
	"strings"
	"testing"
)

func TestDL(t *testing.T) {
	h, err := TORR.Add("magnet:?xt=urn:btih:3522edcc5e979347bf1bc6a99cf12c15b5e66170&tr=http://open.acgtracker.com:1096/announce", "./dl", "mikan@subject-123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(h)
}

func TestDLcompl(t *testing.T) {
	c, err := TORR.DLcompl("3522edcc5e979347bf1bc6a99cf12c15b5e66170")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(c)
}

func TestGetViaPath(t *testing.T) {
	ts, err := TORR.GetViaPath("D:\\mikan-subj\\381666@mikan")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%#+v", ts[0])
}

func TestStr(t *testing.T) {
	str := "jook     lol"
	fs := strings.Fields(str)
	for _, s := range fs {
		fmt.Println(s)
	}
}

func TestAdd(t *testing.T) {
	hash, err := TORR.Add("https://mikanani.me/Download/20221231/3938fb98aa5bbc2d3e0631852accd234feef932e.torrent", "E:\vrhouse", "test123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(hash)
}

func TestCateg(t *testing.T) {
	hs, err := TORR.GetViaCateg("anicat@subj-274234")
	if err != nil {
		t.Error(err)
	}
	t.Log(hs)
}

func TestLs(t *testing.T) {
	list, err := DL.Qbt.TorrentList(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#+v", list)
}
