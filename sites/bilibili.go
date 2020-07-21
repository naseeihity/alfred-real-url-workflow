package sites

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

// BiliID for bilibili method
type BiliID struct {
	RId string
}

// GetRid get real room id
func getRid(rid string) (bool, string, string) {
	const ridURL = "https://api.live.bilibili.com/room/v1/Room/room_init?id=%s"
	var status, id, title = false, "0", ""

	// fetch real room id
	url := fmt.Sprintf(ridURL, rid)
	data, err := httplib.Get(url).String()
	if err != nil {
		return status, id, title
	}

	// conver to json
	res, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return status, id, title
	}

	statusNum, err := res.Get("data").Get("live_status").Int()
	idNum, err := res.Get("data").Get("room_id").Int()
	uid, err := res.Get("data").Get("uid").Int()
	if err != nil {
		return status, id, title
	}

	status = statusNum != 0
	if status {
		title = getBiliRoomName(uid)
	}
	id = strconv.Itoa(idNum)

	log.Println("status:", status, " => id :", id)

	return status, id, title
}

// getBiliRoomName get name of room
func getBiliRoomName(uid int) string {
	const homeURL = "http://api.bilibili.com/x/space/acc/info?mid=%s"
	var id = strconv.Itoa(uid)
	var title = "bilibi_" + id
	url := fmt.Sprintf(homeURL, id)
	data, err := httplib.Get(url).String()
	if err != nil {
		return title
	}

	// conver to json
	res, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return title
	}

	titleStr, err := res.Get("data").Get("name").String()
	if err != nil {
		return title
	}

	title = "bilibili_" + titleStr + "_" + id

	println("the title is:", title)

	return title
}

// GetOneURL get real url of bilibili
func (id BiliID) GetOneURL() (RoomInfo, error) {
	const roomURL = "https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomPlayInfo?room_id=%s&play_url=1&mask=1&qn=1&platform=web"
	var realURL string
	rid := string(id.RId)
	status, rid, title := getRid(rid)

	roomInfo := RoomInfo{
		Title: title,
		URL:   realURL,
	}

	if !status || rid == "0" {
		err := errors.New("Not on live or room not found")
		log.Println("Not on live or room not found")
		return roomInfo, err
	}

	url := fmt.Sprintf(roomURL, rid)

	// conver to json
	res, err := GetJSONRes(url)
	if err != nil {
		log.Println("Bilibili => GetJSONRes Failed:", err)
		return roomInfo, err
	}

	arr, err := res.Get("data").Get("play_url").Get("durl").Array()
	if err != nil {
		log.Println("Not get Url:", err)
		return roomInfo, err
	}
	bestURLInfo, ok := arr[len(arr)-1].(map[string]interface{})
	if !ok {
		err := errors.New("Not get Url")
		log.Println("Not get Url")
		return roomInfo, err
	}
	realURL = bestURLInfo["url"].(string)

	log.Println(title, " ==> ", realURL)
	roomInfo.URL = realURL

	return roomInfo, nil
}

// GetURL used for channel
func (id BiliID) GetURL(ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo, err := id.GetOneURL()
	if err != nil {
		log.Printf("Get bilibili URL of rid-%s Failed:%s", id.RId, err)
	}
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
