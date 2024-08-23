package main

import (
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/subject"
	"github.com/gosuri/uiprogress"
	"time"
)

func main() {
	dler := builtin.NewDownloader("D:\\torrdir\\sqllite_db", false, builtin.NewHttpSeeker())
	fp := subject.FilePath{FileName: &subject.RssFileOpt{Renamed: "testS02E01"}, DirPath: "D:\\torrdir"}

	t, err := dler.Download("https://mikanani.me/Download/20240805/4161be52f9d578a5305188c7c000273171371625.torrent", fp, nil)
	if err != nil {
		panic(err)
	}
	<-t.GotInfo()
	t.DownloadAll()
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
