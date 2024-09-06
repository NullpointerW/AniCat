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
	fp := subject.FilePath{FileName: &subject.RssFileOpt{Renamed: "testS02E01"}, DirPath: "D:\\torrdir"}
     
	t, err := dler.Download("e1cbb423bda0585f978c9b27f6f05263cb590f22", fp, nil )
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
