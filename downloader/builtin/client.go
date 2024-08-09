package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Downloader struct {
	client    *torrent.Client
	nopUpload bool
	baseDir   string
}

func main() {
	cfg := torrent.NewDefaultClientConfig()
	//cfg.DataDir = "db"
	cfg.NoUpload = true
	fop2 := storage.NewFileClientOpts{}
	//fop2.TorrentDirMaker = func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
	//	return "D:\\torrdir"
	//}
	fop2.ClientBaseDir = "torr_db"
	cfg.DefaultStorage = storage.NewFileOpts(fop2)

	// Create a new torrent client with default configuration
	client, err := torrent.NewClient(cfg)

	if err != nil {
		log.Fatalf("error creating torrent client: %v", err)
	}
	defer client.Close()

	// Define the magnet link
	//magnetLink := "magnet:?xt=urn:btih:fbc74348498175b4caec790054e922d28312a106&tr=http%3a%2f%2ft.nyaatracker.com%2fannounce&tr=http%3a%2f%2ftracker.kamigami.org%3a2710%2fannounce&tr=http%3a%2f%2fshare.camoe.cn%3a8080%2fannounce&tr=http%3a%2f%2fopentracker.acgnx.se%2fannounce&tr=http%3a%2f%2fanidex.moe%3a6969%2fannounce&tr=http%3a%2f%2ft.acg.rip%3a6699%2fannounce&tr=https%3a%2f%2ftr.bangumi.moe%3a9696%2fannounce&tr=udp%3a%2f%2ftr.bangumi.moe%3a6969%2fannounce&tr=http%3a%2f%2fopen.acgtracker.com%3a1096%2fannounce&tr=udp%3a%2f%2ftracker.opentrackr.org%3a1337%2fannounce"
	// Add the magnet link to the client
	//t, err := client.AddMagnet(magnetLink)
	trsp := torrent.TorrentSpecFromMetaInfo(read())
	//trsp.Storage = storage.NewFile("D:\\torrdir\\sousou")
	fop := storage.NewFileClientOpts{}
	fop.TorrentDirMaker = func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
		return "D:\\torrdir"
	}
	fop.FilePathMaker = func(opts storage.FilePathMakerOpts) string {
		var parts []string
		fmt.Println("info name=", opts.Info.Name)
		if opts.Info.Name != metainfo.NoName {
			//fmt.Println(opts.Info.Name)
			//rename := "testrename" + filepath.Ext(opts.Info.Name)
			parts = append(parts, opts.Info.Name)
			//parts = append(parts, rename)
		}
		for _, fp := range opts.File.Path {
			fmt.Println(fp)
		}
		fmt.Println(filepath.Join(append(parts, opts.File.Path...)...))
		return filepath.Join(append(parts, opts.File.Path...)...)
	}
	//fop.ClientBaseDir = "D:\\torrdir\\base"
	trsp.Storage = storage.NewFileOpts(fop)
	t, _, _ := client.AddTorrentSpec(trsp)
	if err != nil {
		log.Fatalf("error adding magnet link: %v", err)
	}

	<-t.GotInfo()
	t.DownloadAll()
	go func() {
		for {
			//for _, f := range t.Files() {
			//	fmt.Println(f.Path())
			//	fmt.Println(float64(f.BytesCompleted())/float64(f.Length())*100, "%")
			//}
			fmt.Println(float64(t.BytesCompleted())/float64(t.Info().TotalLength())*100, "%")
			time.Sleep(1 * time.Second)
		}

	}()
	select {
	case <-t.Complete.On():
		log.Println("torrent download complete")
	}
}

func read() *metainfo.MetaInfo {
	resp, _ := http.Get("https://mikanani.me/Download/20240629/799328b2d580a66e25640fdea2d17302501eca08.torrent")
	mi, err := metainfo.Load(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Println(mi)
	return mi
}
