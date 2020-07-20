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
