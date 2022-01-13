package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	new_api = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo?device=phone&platform=h5&scale=3&build=10000&protocol=0,1&format=0,1,2&codec=0,1&room_id="
	uid_api = "http://api.live.bilibili.com/live_user/v1/Master/info?uid="
)

type Live_room_detail struct {
	Code    int    `json: "code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Room_id           int   `json:"room_id"`
		Short_id          int   `json:"short_id"`
		Uid               int   `json:"uid"`
		Need_p2p          int   `json:"need_p2p"`
		Is_hidden         bool  `json:"is_hidden"`
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

type Btuber struct {
	ready       bool
	User_info   Live_user_info
	Room_detail Live_room_detail
}

func New() Btuber {
	var newBtuber Btuber
	return newBtuber
}
func (p *Btuber) Parse_uid(uid string) error {
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
	return nil
}
func (p *Btuber) Get_room_info() error {
	resp, err := http.Get(new_api + strconv.Itoa(p.User_info.Data.Room_id))
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
	target := p.Room_detail.Data.Playurl_info.Playurl.Stream[0].Format[0].Codec[0]
	return target.Url_info[0].Host + target.Base_url + target.Url_info[0].Extra, nil
}
func (p *Btuber) Get_user_name() string {
	return p.User_info.Data.Info.Uname
}
