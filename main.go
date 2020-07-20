package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"realurl/operate"
)

func main() {
	start := time.Now()
	args := os.Args[1:]
	argsLen := len(args)

	switch argsLen {
	case 0:
		operate.PlayFromJSON()
	case 1:
		fmt.Printf("Play directly")
	case 2:
		fmt.Printf("1. read from json 2.Play specific site rooms.")
	case 3:
		fmt.Printf("Add or Play")
	default:
		log.Println("Not valid Arg number")
	}
	log.Printf("Total Time: %.2fs\n", time.Since(start).Seconds())
}
