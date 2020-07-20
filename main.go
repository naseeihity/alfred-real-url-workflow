package main

import (
	"log"
	"os"
	"strings"
	"time"

	"realurl/operate"
)

func main() {
	start := time.Now()
	args := os.Args[1:]
	argsLen := len(args)

	switch argsLen {
	case 0:
		operate.PlayAll()
	case 1:
		// ugly, maybe rewrite in the future
		if strings.ToLower(args[0]) == "play" {
			operate.Play("")
		} else {
			log.Println("Invalid Arg!")
		}
	case 2:
		if strings.ToLower(args[0]) == "play" {
			operate.PlayByPlatform(args[1])
		} else {
			log.Println("Invalid Arg!")
		}
	case 3:
		if strings.ToLower(args[0]) == "play" {
			operate.PlayByID(args[1], args[2])
		} else if strings.ToLower(args[0]) == "add" {
			operate.AddNewRoom(args[1], args[2])
		} else {
			log.Println("Invalid Arg!")
		}
	default:
		log.Println("Not valid Arg number")
	}
	log.Printf("Total Time: %.2fs\n", time.Since(start).Seconds())
}
