package sites

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// YoukuID for youku method
type YoukuID struct {
	RId string
}

// GetOneURL get real url of youku
func (id YoukuID) GetOneURL() (RoomInfo, error) {
	const roomURL = "https://acs.youku.com/h5/mtop.youku.live.com.livefullinfo/1.0/?appKey=24679788"
	const signStr = "%s&%d&24679788&%s"
	const baseURL = "http://lvo-live.youku.com/vod2live/%s_mp4hd2v3.m3u8?&expire=21600&psid=1&ups_ts=%s&vkey="
	rid := string(id.RId)
	title := "youku_" + rid
	roomInfo := RoomInfo{
		Title: title,
		URL:   "",
	}

	// 发送请求获取cookie
	resp, err := http.Get(roomURL)
	if err != nil {
		return roomInfo, err
	}
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		err := errors.New("Get cookies failed")
		return roomInfo, err
	}

	// 构造jar，用于将cookie带到之后的cookie中
	jar, _ := cookiejar.New(nil)
	// URL for cookies to remember. i.e reply when encounter this URL
	cookieURL, _ := url.Parse(roomURL)
	jar.SetCookies(cookieURL, cookies)

	// 请求相关参数处理
	token := string([]byte(cookies[0].Value)[:32])
	timeNow := NowAsUnixMilli()
	dataMap := map[string]string{
		"liveId": rid,
		"app":    "Pc",
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return roomInfo, err
	}
	sign := fmt.Sprintf(signStr, token, timeNow, data)
	sign = GetMD5Hash(sign)
	params := map[string]string{
		"t":    strconv.FormatInt(timeNow, 10),
		"sign": sign,
		"data": string(data),
	}

	res, err := GetJSONResWithCookie(roomURL, params, jar)
	if err != nil {
		return roomInfo, err
	}

	// 处理JSON太麻烦了
	// 这算是个处理样例
	title, err = res.Get("data").Get("data").Get("name").String()
	streamUrls, err := res.Get("data").Get("data").Get("stream").Array()
	if len(streamUrls) == 0 {
		err := errors.New("Get url failed")
		return roomInfo, err
	}
	stream, ok := streamUrls[0].(map[string]interface{})
	if !ok {
		err := errors.New("Not get Url")
		return roomInfo, err
	}
	streamURL := stream["streamName"].(string)
	if err != nil {
		log.Println("Not get url or title:", err)
		return roomInfo, err
	}

	roomInfo.URL = fmt.Sprintf(baseURL, streamURL, strconv.FormatInt(time.Now().UnixNano(), 10))
	roomInfo.Title = title

	return roomInfo, nil
}

// GetURL used for channel
func (id YoukuID) GetURL(ch chan<- RoomInfo, wg *sync.WaitGroup) {
	start := time.Now()
	defer wg.Done()
	roomInfo, err := id.GetOneURL()
	if err != nil {
		log.Println("Get youku URL of rid-%s Failed:%s", id.RId, err)
	}
	ch <- roomInfo
	log.Printf("%.2fs %s\n", time.Since(start).Seconds(), roomInfo.Title)
}
