package main

import (
	"fmt"
	"time"

	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/subject"
	"github.com/gosuri/uiprogress"
)
func main() {
	dler := builtin.NewDownloader("D:\\torrdir\\sqllite_db", false, builtin.NewHttpSeeker())
	go download(dler, "https://mikanani.me/Download/20250105/d7932396b043e69718e5909d81e89c902f409e6d.torrent", "test2.mp4", "D:\\torrdir2")
	go download(dler, "https://mikanani.me/Download/20221225/1aac6b80dd286d6d54eb55adb8d067af535a03c9.torrent", "test.mp4", "D:\\torrdir")
	select {}
}
func download(dw *builtin.Downloader,url string, name string, dir string) {
	fp := subject.FilePath{FileName: &subject.RssFileOpt{Renamed: name}, DirPath: dir}

	t, err := dw.Download(url, fp, nil)
	if err != nil {
		panic(err)
	}
	<-t.GotInfo()
	t.DownloadAll()
	fmt.Println(t.InfoHash().HexString())

	uiprogress.Start()
	bar := uiprogress.AddBar(100)
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "test.mp4"
	})
	bar.AppendCompleted()
	//bar.PrependElapsed()
	go func() {
		for {
			progress := int(float64(t.BytesCompleted()) / float64(t.Length()) * 100)
			err := bar.Set(progress)
			if err != nil {
				panic(err)
			}
			time.Sleep(40 * time.Millisecond)
		}
	}()
	<-t.Complete.On()
}
