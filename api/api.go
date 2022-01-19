package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	room_api = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo?platform=web&protocol=0,1&format=0,1,2&codec=0,1&qn=0&room_id="
	uid_api  = "http://api.live.bilibili.com/live_user/v1/Master/info?uid="
)

type Live_room_detail struct {
	// This struct is directly derived from the JSON response of room_api
	Code    int    `json: "code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Room_id           int   `json:"room_id"`
		Short_id          int   `json:"short_id"`
		Uid               int   `json:"uid"`
		Need_p2p          int   `json:"need_p2p"`
		Ishidden          bool  `json:"is_hidden"`
		Is_locked         bool  `json:"is_locked"`
		Is_portrait       bool  `json:"is_portrait"`
		Live_status       int   `json:"live_status"`
		Hidden_till       int   `json:"hidden_till"`
		Lock_till         int   `json:"lock_till"`
		Encrypted         bool  `json:"encrypted"`
		Pwd_verified      bool  `json:"pwd_verified"`
		Live_time         int   `json:"live_time"`
		Room_shield       int   `json:"room_shield"`
		All_special_types []int `json:"all_special_types"`
		Playurl_info      struct {
			Conf_json string `json:"conf_json"`
			Playurl   struct {
				Cid       int `json:"cid"`
				G_qn_desc []struct {
					Qn       int    `json:"qn"`
					Desc     string `json:"desc"`
					Hdr_desc string `json:"hdr_desc"`
				} `json:"g_qn_desc"`
				Stream []struct {
					Protocol_name string `json:"protocol_name"`
					Format        []struct {
						Format_name string `json:"format_name"`
						Codec       []struct {
							Codec_name string `json:"codec_name"`
							Current_qn int    `json:"current_qn"`
							Accept_qn  []int  `json:"accept_qn"`
							Base_url   string `json:"base_url"`
							Url_info   []struct {
								Host       string `json:"host"`
								Extra      string `json:"extra"`
								Stream_ttl int    `json:"stream_ttl"`
							} `json:"url_info"`
							Hdr_qn int `json:"hdr_qn"`
						} `json:"codec"`
					} `json:"format"`
				} `json:"stream"`
				P2p_data struct {
					P2p       bool `json:"p2p"`
					P2p_type  int  `json:"p2p_type"`
					M_p2p     bool `json:"M_p2p"`
					M_servers int  `json:"m_servers"`
				} `json:"p2p_data"`
				Dolby_qn int `json:"dolby_qn"`
			} `json:"playurl"`
		} `json:"playurl_info"`
	} `json:"data"`
}
type Live_user_info struct {
	// The struct is directly derived from the JSON response of uid_api
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			Uid             int    `json:"uid"`
			Uname           string `json:"uname"`
			Face            string `json:"face"`
			Official_verify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Gender int `json:"gender"`
		} `json:"info"`
		Exp struct {
			Master_level struct {
				Level   int   `json:"level"`
				Color   int   `json:"color"`
				Current []int `json:"current"`
				Next    []int `json:"next"`
			} `json:"master_level"`
		} `json:"exp"`
		Follower_num   int    `json:"follower_num"`
		Room_id        int    `json:"room_id"`
		Medal_name     string `json:"medal_name"`
		Glory_count    int    `json:"glory_count"`
		Pendant        string `json:"pendant"`
		Link_group_num int    `json:"link_group_num"`
		Room_news      struct {
			Content    string `json:"content"`
			Ctime      string `json:"ctime"`
			Ctime_text string `json:"ctime_text"`
		} `json:"room_news"`
	} `json:"data"`
}
type Stream_url struct {
	Url      string
	Qn       int
	Codec    string
	Format   string
	Protocol string
}

func (p *Btuber) Show_avaible_urls() []Stream_url {
	var urls []Stream_url
	for _, protocal := range p.Room_detail.Data.Playurl_info.Playurl.Stream {
		for _, format := range protocal.Format {
			for _, codec := range format.Codec {
				if len(codec.Url_info) > 0 {
					urls = append(urls, Stream_url{Url: codec.Url_info[0].Host + codec.Base_url + codec.Url_info[0].Extra, Qn: codec.Current_qn, Codec: codec.Codec_name, Format: format.Format_name, Protocol: protocal.Protocol_name})
				}
			}
		}
	}
	return urls
}
func (p *Btuber) Get_room_news() (string, error) {
	p.ready = false
	p.Parse_uid(strconv.Itoa(p.User_info.Data.Info.Uid))
	if !p.ready {
		return "", errors.New("Btuber not initialized with legal uid")
	}
	news := p.User_info.Data.Info.Uname + " posted:\n"
	news += p.User_info.Data.Room_news.Content + "\n on " + p.User_info.Data.Room_news.Ctime_text
	return news, nil
}

type Btuber struct {
	ready       bool
	User_info   Live_user_info   // Basic user info | possibly won't change | filled when initialize
	Room_detail Live_room_detail // Live room info | could change any time | fetched when used
}

func New() Btuber {
	// Ctor
	var newBtuber Btuber
	newBtuber.ready = false // the user is forced to initialize with uid
	return newBtuber
}
func (p *Btuber) On_live() (bool, error) {
	// Return if the room is currently in Livestream
	// Must called after Parse_uid
	if !p.ready {
		return false, errors.New("Btuber not initialized with legal uid")
	}
	err := p.Get_room_info() // Refreshing the room info first
	if err != nil {
		return false, err
	}
	return p.Room_detail.Data.Live_status == 1, nil
}
func (p *Btuber) Subscribe() (chan int, error) {
	// return a notify channel subscibing to the room
	// notify would recieve a value when the room is live
	if !p.ready {
		return nil, errors.New("Btuber not initialized with legal uid")
	}
	notify := make(chan int)
	go p.Test_on_live(notify)
	return notify, nil
}
func (p *Btuber) Test_on_live(notify chan int) {
	// Repeatedly check whether the room is live
	// send 1 to notify channel when live
	liveness := false
	t := time.NewTicker(5000 * time.Millisecond) // Politeness
	defer t.Stop()
	live, _ := p.On_live()
	if !liveness && live {
		notify <- 1
	}
	liveness = live
	for {
		select {
		case <-t.C:
			live, _ := p.On_live()
			if !liveness && live {
				notify <- 1
			}
			liveness = live
		}
	}
}
func (p *Btuber) Parse_uid(uid string) error {
	// Fill a Btuber instance with information fetched from his/her uid
	// Return an error if failed in http or wrong uid
	resp, err := http.Get(uid_api + uid)
	if err != nil {
		fmt.Println("Error in parsing UID")
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &p.User_info)
	if p.User_info.Data.Info.Uid == 0 {
		return errors.New("Null User")
	}
	p.ready = true
	return nil
}
func (p *Btuber) Get_room_info() error {
	// Fill a Btuber instance with room information fetched from its room_id
	// Must called after Parse_uid
	if !p.ready {
		return errors.New("Btuber not initialized with legal uid")
	}
	resp, err := http.Get(room_api + strconv.Itoa(p.User_info.Data.Room_id))
	if err != nil {
		fmt.Println("Error in parsing Roomid")
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &p.Room_detail)
	if p.Room_detail.Data.Room_id == 0 {
		return errors.New("Null Room ID")
	}
	return nil
}
func (p *Btuber) Get_url() (string, error) {
	// Return the default livestream url of the Btuber
	// Must called after Get_room_info
	if !p.ready {
		return "", errors.New("Btuber not initialized with legal uid")
	}
	liveness, _ := p.On_live()
	if !liveness {
		return "", errors.New("Room not in Live")
	}
	target := p.Room_detail.Data.Playurl_info.Playurl.Stream[0].Format[0].Codec[0]
	return target.Url_info[0].Host + target.Base_url + target.Url_info[0].Extra, nil
}
func (p *Btuber) Get_user_name() (string, error) {
	// Return User name of the Btuber
	if !p.ready {
		return "", errors.New("Btuber not initialized with legal uid")
	}
	return p.User_info.Data.Info.Uname, nil
}
