package sites

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

// Platform interface wiht GetURL method
type Platform interface {
	GetURL(chan<- RoomInfo, *sync.WaitGroup)
	GetOneURL() (RoomInfo, error)
}

//RoomInfo url and title
type RoomInfo struct {
	URL   string
	Title string
}

//GetJSONRes get and convert to json
func GetJSONRes(url string) (*simplejson.Json, error) {
	data, err := httplib.Get(url).String()
	if err != nil {
		log.Fatal("http request error:", err)
		return nil, err
	}

	// conver to json
	res, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Fatal("json convert error:", err)
		return nil, err
	}

	return res, nil
}

// GetWithHead a content-type head
func GetWithHead(url string, header map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)

	return resp, err
}

//GetJSONResWithCookie 设置了cookie和jar
func GetJSONResWithCookie(url string, params map[string]string, jar *cookiejar.Jar) (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置Get请求的params
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// 响应体转为string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		data := string(bodyBytes)

		// conver to json
		res, err := simplejson.NewJson([]byte(data))
		if err != nil {
			log.Fatal("json convert error:", err)
			return nil, err
		}

		err = errors.New("Get request failed")

		return res, nil
	}

	return nil, err
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//NowAsUnixMilli  current time in seconds since the Epoch
func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//GetMD5Hash hast a string using md5
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
