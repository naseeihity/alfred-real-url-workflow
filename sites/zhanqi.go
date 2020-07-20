package sites

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ZhanqiID for zhanqi method
type ZhanqiID struct {
	RId string
}

// GetOneURL get real url of zhanqi
func (id ZhanqiID) GetOneURL() (RoomInfo, error) {
	const roomJSON = "https://m.zhanqi.tv/api/static/v2.1/room/domain/%s.json"
	const roomURL = "https://dlhdl-cdn.zhanqi.tv/zqlive/%s.flv"
	rid := string(id.RId)
	title := "zhanqi_" + rid
	roomInfo := RoomInfo{
		Title: title,
		URL:   "",
	}

	url := fmt.Sprintf(roomJSON, rid)
	res, err := GetJSONRes(url)
	if err != nil {
		log.Println("Zhanqi => GetJSONRes Failed:", err)
		return roomInfo, err
	}

	videoID, err := res.Get("data").Get("videoId").String()
	if err != nil {
		log.Println("Zhanqi => Get VideoId Failed:", err)
		return roomInfo, err
	}

	if len(videoID) != 0 {
		roomInfo.URL = fmt.Sprintf(roomURL, videoID)
	}
	return roomInfo, nil
}

// GetURL used for channel
func (id ZhanqiID) GetURL(ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo, err := id.GetOneURL()
	if err != nil {
		log.Fatalf("Get zhanqi URL of rid-%s Failed:%s", id.RId, err)
	}
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
