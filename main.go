package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"./sites"
)

const (
	playList = "./playlist.m3u8"
)

func main() {
	var roomInfos []sites.RoomInfo
	var wg sync.WaitGroup

	start := time.Now()
	rooms := []string{"21753173", "888", "41515"}
	ch := make(chan sites.RoomInfo)

	// TODO: read cmd from Arg

	// get all roominfos concurency
	for _, room := range rooms {
		wg.Add(1)
		go sites.GetBilibiliURL(room, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for room := range ch {
		roomInfos = append(roomInfos, room)
	}

	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

	// new or edit file and open
	openPlayList(roomInfos)
}

func openPlayList(roomInfos []sites.RoomInfo) {
	f, err := os.OpenFile(playList, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, info := range roomInfos {
			roomTitle := fmt.Sprintf("#EXTINF:-1,%s\n", info.Title)
			_, err = f.Write([]byte(roomTitle))
			roomURL := fmt.Sprintf("%s\n", info.URL)
			_, err = f.Write([]byte(roomURL))
		}
		cmd := exec.Command("open", playList)
		cmd.Start()
	}
}
