package main

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"os"
	"strconv"
	"time"
	"zhujiaxu.com/bilibili-livestream-recorder/api"
)

func main() {
    // Basic usage example
    // Calling with go run . [uid]
    // For example go run . 123123123123 or bilibili-livestream-recorder 123123123123
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No Arguments")
		return
	}
	uid := args[1]
	btuber := api.New()
	btuber.Parse_uid(uid)
	btuber.Get_room_info()
	monitor(&btuber)
	return
}

type video_slice struct {
	url string
	id  int
}

func downloader(input chan video_slice, output chan video_slice, done chan int) error {
	var to_download video_slice
	client := grab.NewClient()
	for {
		select {
		case to_download = <-input:
			req, _ := grab.NewRequest(".", to_download.url)
			req.Filename = strconv.Itoa(to_download.id) + ".flv"
			fmt.Println("Downloading Slice", req.Filename)
			resp := client.Do(req)
			t := time.NewTicker(500 * time.Millisecond)
			defer t.Stop()

		Loop:
			for {
				select {
				case <-t.C:
					if resp.BytesComplete() < 10000000 {
						fmt.Printf("  transferred %.2f KBs\n", float64(resp.BytesComplete())/1000.0)
					} else {
						fmt.Printf("  transferred %.2f MBs\n", float64(resp.BytesComplete())/1000000.0)
					}

				case <-resp.Done:
					// download is complete
					break Loop
				}
			}
			output <- to_download
			break
		case <-done:
			return nil
		}
	}
}

func /*(p *api.Btuber)*/ monitor(p *api.Btuber) error {
    // Monitoring the Btuber
    // Download video slice when room is live
	queue := make(chan video_slice)
	output := make(chan video_slice)
	done := make(chan int)
	notify, err := p.Subscribe()
	if err != nil {
		return err
	}
	go downloader(queue, output, done)
	count := 0
	name, _ := p.Get_user_name()
	fmt.Printf("Monitoring %v's Room\n", name)
	for {
		<-notify

		url, err := p.Get_url()
		if err != nil {
			done <- 1
			return err
		}
		count += 1
		fmt.Printf("Starting download slice %d\n", count)
		queue <- video_slice{url: url, id: count}
		<-output
	}
	done <- 1
	return nil
}
