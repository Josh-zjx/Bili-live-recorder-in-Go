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
	queue := make(chan video_slice)
	output := make(chan video_slice)
	done := make(chan int)
	go downloader(queue, output, done)
	count := 0
	prev, err := p.Get_url()
	if err != nil {
		done <- 1
		return err
	}
	queue <- video_slice{url: prev, id: count}
	var url string
	for {
		<-output
		url, err = p.Get_url()
		if err != nil {
			done <- 1
			return err
		}
		if url != prev {
			count += 1
			queue <- video_slice{url: url, id: count}
			prev = url
		}
	}
	done <- 1
	return nil
}
