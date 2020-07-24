package operate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"realurl/sites"
	"strings"
	"sync"
)

const (
	ridJSON   = "roomlist.json"
	playList  = "playlist.m3u8"
	playListT = "_playlist.m3u8"
)

// getRoomsFromJSON get room rid from cahced json file
// This function shows how to convert a json file into a golang map
// as well as convert []interface{} to []string
func getRoomsFromJSON(p string) map[string][]string {
	var rooms map[string][]interface{}
	// should not int the map as a nil map which
	// will cause assignment to entry in nil map panic when do roomMap[key] = value
	roomMap := make(map[string][]string)

	// read from json file and unmarshal
	f, err := ioutil.ReadFile(getPath(ridJSON))
	if err != nil {
		log.Fatal("open file err:", err)
	}
	json.Unmarshal(f, &rooms)

	// convert map interface to map string
	for platform, roomInfo := range rooms {
		if len(p) != 0 && p != platform {
			continue
		}
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

	err := setM3U8File(roomInfos, false)

	return err
}

func setM3U8File(roomInfos []sites.RoomInfo, temporary bool) error {
	var buffer bytes.Buffer
	for _, info := range roomInfos {
		roomTitle := fmt.Sprintf("#EXTINF:-1,%s\n", info.Title)
		roomURL := fmt.Sprintf("%s\n", info.URL)
		// using buffer to do string concat
		buffer.WriteString(roomTitle + roomURL)
	}
	txt := buffer.String()

	if temporary {
		err := ioutil.WriteFile(getPath(playListT), []byte(txt), 0666)
		return err
	}
	err := ioutil.WriteFile(getPath(playList), []byte(txt), 0666)
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
	switch strings.ToLower(platform) {
	case "douyu":
		for _, id := range rids {
			roomIds = append(roomIds, sites.DouyuID{RId: id})
		}
	case "bilibili":
		for _, id := range rids {
			roomIds = append(roomIds, sites.BiliID{RId: id})
		}
	case "zhanqi":
		for _, id := range rids {
			roomIds = append(roomIds, sites.ZhanqiID{RId: id})
		}
	case "youku":
		for _, id := range rids {
			roomIds = append(roomIds, sites.YoukuID{RId: id})
		}
	case "huya":
		for _, id := range rids {
			roomIds = append(roomIds, sites.HuyaID{RId: id})
		}
	default:
		log.Printf("No such Platform")
	}

	return roomIds
}

// playFromJSON conver json to playlist and play
func playFromJSON(p string) error {
	roomMap := getRoomsFromJSON(p)

	if len(roomMap) == 0 {
		err := errors.New("No valid room rid")
		return err
	}

	if err := convertToM3U8(roomMap); err != nil {
		log.Fatal("Convert json to M3U8 error:", err)
		return err
	}

	Play(playList)

	return nil
}

// PlayAll all room in local json
func PlayAll() {
	playFromJSON("")
}

// PlayByPlatform filter rooms by platform
func PlayByPlatform(p string) {
	playFromJSON(p)
}

// PlayByID play one room
func PlayByID(p string, rid string) {
	rids := []string{rid}
	roomIds := getPlatRids(p, rids)

	for _, rid := range roomIds {
		roomInfo, err := rid.GetOneURL()
		if err != nil {
			log.Println("Get room url by one id failed: ", err)
		}
		roomInfos := []sites.RoomInfo{roomInfo}
		setM3U8File(roomInfos, true)
		Play(playListT)
	}
}

//AddNewRoom add new room to local json
func AddNewRoom(p string, rids string) {
	roomMap := getRoomsFromJSON("")
	ridSlice := strings.Split(rids, ",")
	var ridSet []string

	// deduplication
	for _, rid := range ridSlice {
		rid := strings.TrimSpace(rid)
		if _, ok := sites.Find(roomMap[p], rid); ok {
			continue
		}
		ridSet = append(ridSet, rid)
	}
	if len(ridSet) == 0 {
		log.Printf("Room id(s) already exist")
		return
	}

	// update
	roomMap[p] = append(roomMap[p], ridSet...)
	data, err := json.MarshalIndent(roomMap, "", "    ")
	if err != nil {
		log.Fatal("add new room faild when covert map to json:", err)
	} else {
		err := ioutil.WriteFile(getPath(ridJSON), data, 0666)
		if err != nil {
			log.Fatal("add new room faild when write to json file:", err)
			return
		}
		log.Printf("Add new room success!")
	}
}

// Play from playList directly
func Play(f string) {
	if len(f) == 0 {
		f = playList
	}
	cmd := exec.Command("open", getPath(f))
	cmd.Start()
}

func getPath(p string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(dir, p)
}
