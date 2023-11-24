package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

func main() {
	// Create a new torrent client with default configuration
	client, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatalf("error creating torrent client: %v", err)
	}
	defer client.Close()

	// Define the magnet link
	magnetLink := "magnet:?xt=urn:btih:fbc74348498175b4caec790054e922d28312a106&tr=http%3a%2f%2ft.nyaatracker.com%2fannounce&tr=http%3a%2f%2ftracker.kamigami.org%3a2710%2fannounce&tr=http%3a%2f%2fshare.camoe.cn%3a8080%2fannounce&tr=http%3a%2f%2fopentracker.acgnx.se%2fannounce&tr=http%3a%2f%2fanidex.moe%3a6969%2fannounce&tr=http%3a%2f%2ft.acg.rip%3a6699%2fannounce&tr=https%3a%2f%2ftr.bangumi.moe%3a9696%2fannounce&tr=udp%3a%2f%2ftr.bangumi.moe%3a6969%2fannounce&tr=http%3a%2f%2fopen.acgtracker.com%3a1096%2fannounce&tr=udp%3a%2f%2ftracker.opentrackr.org%3a1337%2fannounce"

	// Add the magnet link to the client
	t, err := client.AddMagnet(magnetLink)
	if err != nil {
		log.Fatalf("error adding magnet link: %v", err)
	}

	// Wait for the torrent information to be fetched or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	<-t.GotInfo()
	t.DownloadAll()

	go func() {
		for {
			fmt.Printf(" w %d", t.BytesCompleted())
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-t.Complete.On():
		log.Println("torrent download complete")
	case <-ctx.Done():
		log.Println("torrent download timeout")
	}

	// Stream the torrent content to stdout or handle it as needed
	//reader := t.NewReader()
	//io.Copy(os.Stdout, reader)
}
