package sites

import (
	"log"
	"sync"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

// Platform interface wiht GetURL method
type Platform interface {
	GetURL(chan<- RoomInfo, *sync.WaitGroup)
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
