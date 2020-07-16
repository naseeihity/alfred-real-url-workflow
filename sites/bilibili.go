package sites

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

//RoomInfo url and title
type RoomInfo struct {
	URL   string
	Title string
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

	title = getRoomName(uid)

	status = statusNum != 0
	id = strconv.Itoa(idNum)

	log.Println("status:", status, " => id :", id)

	return status, id, title
}

// GetRoomName get name of room
func getRoomName(uid int) string {
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

// GetOneBilibiliURL get real url of bilibili
func GetOneBilibiliURL(rid string) RoomInfo {
	const roomURL = "https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomPlayInfo?room_id=%s&play_url=1&mask=1&qn=0&platform=web"
	var realURL = ""
	status, id, title := getRid(rid)

	var roomInfo RoomInfo
	roomInfo.Title = title
	roomInfo.URL = realURL

	if !status || id == "0" {
		log.Println("Not on live or room not found")
		return roomInfo
	}

	url := fmt.Sprintf(roomURL, id)
	data, err := httplib.Get(url).String()
	if err != nil {
		log.Println("http request error:", err)
		return roomInfo
	}

	// conver to json
	res, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Println("json convert error:", err)
		return roomInfo
	}

	arr, err := res.Get("data").Get("play_url").Get("durl").Array()
	if err != nil {
		log.Println("Not get Url:", err)
		return roomInfo
	}
	bestURLInfo, ok := arr[len(arr)-1].(map[string]interface{})
	if !ok {
		log.Println("Not get Url")
		return roomInfo
	}
	realURL = bestURLInfo["url"].(string)

	log.Println(title, " ==> ", realURL)
	roomInfo.URL = realURL

	return roomInfo
}

// GetBilibiliURL used for channel
func GetBilibiliURL(rid string, ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo := GetOneBilibiliURL(rid)
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
