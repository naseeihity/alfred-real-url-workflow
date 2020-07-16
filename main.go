package main

import (
	"log"
	"sync"
	"time"

	"./sites"
)

func main() {
	start := time.Now()
	var rooms []string
	var roomInfos []sites.RoomInfo
	rooms = []string{"364715", "21180272", "2656132"}

	ch := make(chan sites.RoomInfo)
	var wg sync.WaitGroup

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

	for i, info := range roomInfos {
		log.Printf("%d ==> %s: %s\n", i+1, info.Title, info.URL)
	}

	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
