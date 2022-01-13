package api

import (
	"testing"
)

func Test_New_Btuber(t *testing.T) {
	user := New()
	err := user.Parse_uid("271815944") // Correct Input
	if err != nil {
		t.Errorf("Parse UID error")
	}
	err = user.Get_room_info()
	if err != nil {
		t.Errorf("Parse room info error")
	}
	t.Log(user)
	t.Logf(user.Get_user_name())
}

func Test_Parse_uid_wrong(t *testing.T) {
	user := New()
	err := user.Parse_uid("asdA") // Wrong Input
	if err == nil {
		t.Errorf("Parse UID wrong")
	}
}

func Test_Get_Stream_info(t *testing.T) {
	user := New()
	err := user.Parse_uid("271815944")
	if err != nil {
		t.Errorf("Parse UID failed")
	}
	err = user.Get_room_info()
	if err != nil {
		t.Errorf("Parse Room failed")
	}
	url, _ := user.Get_url()
	t.Logf("%v", url)
	t.Logf("roomid %v", user.Room_detail.Data.Room_id)
}

func Test_close_clip(t *testing.T) {
	user := New()
	_ = user.Parse_uid("105022844")
	_ = user.Get_room_info()
	url, _ := user.Get_url()
	t.Logf("%v", url)
	url, _ = user.Get_url()
	t.Logf("%v", url)
}
