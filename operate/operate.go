package operate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"realurl/sites"
	"sync"
)

const (
	ridJSON  = "roomlist.json"
	playList = "playlist.m3u8"
)

// getRoomsFromJSON get room rid from cahced json file
// This function shows how to convert a json file into a golang map
// as well as convert []interface{} to []string
func getRoomsFromJSON() map[string][]string {
	var rooms map[string][]interface{}
	// should not int the map as a nil map which
	// will cause assignment to entry in nil map panic when do roomMap[key] = value
	roomMap := make(map[string][]string)

	// read from json file and unmarshal
	f, err := ioutil.ReadFile(ridJSON)
	if err != nil {
		log.Fatal("open file err:", err)
	}
	json.Unmarshal(f, &rooms)

	// convert map interface to map string
	for platform, roomInfo := range rooms {
		var roomRids []string
		for _, info := range roomInfo {
			roomRids = append(roomRids, info.(string))
		}
		roomMap[platform] = roomRids
	}

	return roomMap
}

func convertToM3U8(rooms map[string][]string) error {
	var roomInfos []sites.RoomInfo
	var wg sync.WaitGroup
	ch := make(chan sites.RoomInfo)

	for platform, rids := range rooms {
		wg.Add(1)
		go getRoomInfoByPlatform(platform, rids, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for room := range ch {
		roomInfos = append(roomInfos, room)
	}

	err := setM3U8File(roomInfos)

	return err
}

func setM3U8File(roomInfos []sites.RoomInfo) error {
	var buffer bytes.Buffer
	for _, info := range roomInfos {
		roomTitle := fmt.Sprintf("#EXTINF:-1,%s\n", info.Title)
		roomURL := fmt.Sprintf("%s\n", info.URL)
		// using buffer to do string concat
		buffer.WriteString(roomTitle + roomURL)
	}
	txt := buffer.String()

	err := ioutil.WriteFile(playList, []byte(txt), 0666)
	return err
}

func getRoomInfoByPlatform(platform string, rids []string, ch chan<- sites.RoomInfo, wg *sync.WaitGroup) {
	var wg2 sync.WaitGroup
	ch2 := make(chan sites.RoomInfo)
	defer wg.Done()

	roomIds := getPlatRids(platform, rids)

	for _, rid := range roomIds {
		wg2.Add(1)
		go rid.GetURL(ch2, &wg2)
	}

	go func() {
		wg2.Wait()
		close(ch2)
	}()

	for room := range ch2 {
		ch <- room
	}
}

func getPlatRids(platform string, rids []string) []sites.Platform {
	var roomIds []sites.Platform
	// ugly, maybe rewrite in the future
	switch platform {
	case "bilibili":
		for _, id := range rids {
			roomIds = append(roomIds, sites.BiliID{RId: id})
		}
	case "zhanqi":
		for _, id := range rids {
			roomIds = append(roomIds, sites.ZhanqiID{RId: id})
		}
	}
	return roomIds
}

// PlayFromJSON conver json to playlist and play
func PlayFromJSON() error {
	roomMap := getRoomsFromJSON()

	if len(roomMap) == 0 {
		err := errors.New("No valid room rid")
		return err
	}

	if err := convertToM3U8(roomMap); err != nil {
		log.Fatal("Convert json to M3U8 error:", err)
		return err
	}

	Play()

	return nil
}

// Play from playList directly
func Play() {
	cmd := exec.Command("open", playList)
	cmd.Start()
}
