package main

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/jessevdk/go-flags"
	"os"
	"strconv"
	"time"
	"zhujiaxu.com/bilibili-livestream-recorder/api"
)

type Option struct {
	Uid       string `short:"u" long:"uid" required:"true" description:"target uid"`
	Name      bool   `short:"n" long:"name"  description:"query username"`
	Listen    bool   `short:"l" long:"listen"  description:"listen to the room"`
	Show_urls bool   `short:"s" long:"show"  description:"show all urls"`
}

func main() {
	// Basic usage example
	// Calling with go run . [uid]
	// For example go run . 123123123123
	// or bilibili-livestream-recorder 123123123123
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No Arguments")
		return
	}
	var opt Option
	flags.Parse(&opt)
	if opt.Name {
		Get_name(opt.Uid)
	}
	if opt.Show_urls {
		Show_urls(opt.Uid)
	}
	if opt.Listen {
		Listen(opt.Uid)
	}
	//uid := args[1]
	//btuber := api.New()
	//btuber.Parse_uid(uid)
	//btuber.Get_room_info()
	//monitor(&btuber)
	return
}

type video_slice struct {
	url      string
	filename string
}

func Show_urls(uid string) {
	btuber := api.New()
	err := btuber.Parse_uid(uid)
	if err != nil {
		return
	}
	btuber.Get_room_info()
	urls := btuber.Show_avaible_urls()
	for _, url := range urls {
		fmt.Printf("Protocol %v, Format %v, Codec %v, Qn %d, URL: %v\n", url.Protocol, url.Format, url.Codec, url.Qn, url.Url)
	}

}
func Listen(uid string) {
	btuber := api.New()
	err := btuber.Parse_uid(uid)
	if err != nil {
		return
	}
	btuber.Get_room_info()
	monitor(&btuber)
	return
}
func Get_name(uid string) {
	btuber := api.New()
	btuber.Parse_uid(uid)
	name, err := btuber.Get_user_name()
	if err != nil {
		fmt.Printf("Error in getting name for uid %v\n", uid)
		return
	}
	fmt.Printf("%v\n", name)
	return
}

func downloader(input chan video_slice, output chan video_slice, done chan int) error {
	var to_download video_slice
	client := grab.NewClient()
	for {
		select {
		case to_download = <-input:
			req, _ := grab.NewRequest(".", to_download.url)
			req.Filename = to_download.filename + ".flv"
			fmt.Println("Downloading Slice", req.Filename)
			fmt.Printf("Download will start soon\n")
			resp := client.Do(req)
			fmt.Println("0 transferred")
			t := time.NewTicker(500 * time.Millisecond)
			defer t.Stop()

		Loop:
			for {
				select {
				case <-t.C:
					fmt.Printf("\x1b[A")
					if false && resp.BytesComplete() < 10000000 {
						fmt.Printf("  transferred %.2f KBs\n", float64(resp.BytesComplete())/1000.0)
					} else {
						fmt.Printf("  transferred %.2f MBs\n", float64(resp.BytesComplete())/1000000.0)
					}

				case <-resp.Done:
					fmt.Println("Slice Download finished")
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

func monitor(p *api.Btuber) error {
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
		url, _ := p.Get_url()
		count += 1
		fmt.Printf("Starting download slice %d\n", count)
		queue <- video_slice{url: url, filename: p.User_info.Data.Info.Uname + "-" + strconv.Itoa(count)}
		<-output
		fmt.Println("Got one!")
	}
	done <- 1
	return nil
}
