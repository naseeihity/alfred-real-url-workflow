package sites

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// HuyaID for huya method
type HuyaID struct {
	RId string
}

// GetOneURL get real url of huya
func (id HuyaID) GetOneURL() (RoomInfo, error) {
	const roomURL = "https://m.huya.com/%s"
	rid := string(id.RId)
	title := "huya_" + rid
	roomInfo := RoomInfo{
		Title: title,
		URL:   "",
	}

	url := fmt.Sprintf(roomURL, rid)
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36",
	}
	resp, err := GetWithHead(url, header)
	if err != nil {
		return roomInfo, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		regStr := `liveLineUrl = "([\s\S]*?)";`
		regTitle := `<title>([\s\S]*?)_`
		re := regexp.MustCompile(regStr)
		reT := regexp.MustCompile(regTitle)
		partArr := re.FindAllSubmatch([]byte(bodyString), -1)
		tArr := reT.FindAllSubmatch([]byte(bodyString), -1)
		if len(tArr) != 0 {
			title = title + "_" + string(tArr[0][1])
		}

		if len(partArr) != 0 {
			liveLineURL := string(partArr[0][1])
			log.Println(liveLineURL)
			if strings.Contains(liveLineURL, "replay") {
				roomInfo.URL = "https:" + liveLineURL
				roomInfo.Title = title + "_replay"
			} else {
				re := regexp.MustCompile(`_\d{4}.m3u8`)
				liveLineURL = string(re.ReplaceAll([]byte(liveLineURL), []byte(".m3u8")))
				roomInfo.URL = "https:" + liveLineURL
				roomInfo.Title = title
			}
		}
	}

	return roomInfo, nil
}

// GetURL used for channel
func (id HuyaID) GetURL(ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo, err := id.GetOneURL()
	if err != nil {
		log.Fatalf("Get huya URL of rid-%s Failed:%s", id.RId, err)
	}
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
