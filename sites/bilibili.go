package sites

import (
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

const ()

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

// GetBilibiliURL get real url of bilibili
func GetBilibiliURL(rid string) (string, string) {
	const roomURL = "https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomPlayInfo?room_id=%s&play_url=1&mask=1&qn=0&platform=web"
	var realURL = ""
	status, id, title := getRid(rid)

	if !status || id == "0" {
		log.Println("Not on live or room not found")
		return realURL, title
	}

	url := fmt.Sprintf(roomURL, id)
	data, err := httplib.Get(url).String()
	if err != nil {
		log.Println("http request error:", err)
		return realURL, title
	}

	// conver to json
	res, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Println("json convert error:", err)
		return realURL, title
	}

	arr, err := res.Get("data").Get("play_url").Get("durl").Array()
	if err != nil {
		log.Println("Not get Url:", err)
		return realURL, title
	}
	bestURLInfo, ok := arr[len(arr)-1].(map[string]interface{})
	if !ok {
		log.Println("Not get Url")
		return realURL, title
	}
	realURL = bestURLInfo["url"].(string)

	log.Println(title, " ==> ", realURL)

	return realURL, title
}
