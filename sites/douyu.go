package sites

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

// DouyuID for zhanqi method
type DouyuID struct {
	RId string
}

func getDouyuRoomName(rid string) (string, error) {
	const homeURL = "http://open.douyucdn.cn/api/RoomApi/room/%s"

	var title = "douyu_%s_%s_" + rid
	url := fmt.Sprintf(homeURL, rid)

	res, err := GetJSONRes(url)
	if err != nil {
		log.Println("Douyu => getDouyuRoomName Failed:", err)
		return "", err
	}

	if code, err := res.Get("error").Int(); err != nil || code != 0 {
		err := errors.New("Get room info failed")
		return "", err
	}

	roomName, err := res.Get("data").Get("room_name").String()
	ownerName, err := res.Get("data").Get("owner_name").String()
	if err != nil {
		log.Println("Not get roomName or ownerName:", err)
		return "", err
	}

	title = fmt.Sprintf(title, roomName, ownerName)
	return title, nil
}

// GetOneURL get real url of zhanqi
// TODO: Another method, this may miss some rooms
func (id DouyuID) GetOneURL() (RoomInfo, error) {
	const roomURL = "https://www.douyu.com/%s"
	const URL = "http://tx2play1.douyucdn.cn/live/%s.flv?uuid="
	rid := string(id.RId)
	title := "douyu_" + rid
	roomInfo := RoomInfo{
		Title: title,
		URL:   "",
	}

	url := fmt.Sprintf(roomURL, rid)
	resp, err := GetWithHead(url, nil)
	if err != nil {
		log.Fatal("Get douyu url failed:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		regStr := `dyliveflv1a*/([\d\w]*?)_`
		re := regexp.MustCompile(regStr)
		// 注意这里如果要匹配正则表达式括号里的值需要用submatch
		partArr := re.FindAllSubmatch([]byte(bodyString), -1)
		if len(partArr) != 0 {
			part := partArr[0][1]
			liveURL := fmt.Sprintf(URL, string(part))
			roomInfo.URL = liveURL
		}
	}

	if len(roomInfo.URL) > 0 {
		title, err = getDouyuRoomName(rid)
		if err != nil {
			log.Println("Get douyu room name failed:", err)
		} else {
			roomInfo.Title = title
		}
	}

	return roomInfo, nil
}

// GetURL used for channel
func (id DouyuID) GetURL(ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo, err := id.GetOneURL()
	if err != nil {
		log.Fatalf("Get Douyu URL of rid-%s Failed:%s", id.RId, err)
	}
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
