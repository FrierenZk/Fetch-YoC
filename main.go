package main

import (
	. "./Debug"
	"./Fetch"
	"fmt"
	"log"
	"time"
)

const checkFreq = time.Minute * 15

//var timeRecord time.Time

func main() {
	err := LogInit()
	if err != nil {
		log.Fatal(err)
	}
	defer LogExit(nil)
	var update = make(chan string, 1)
	update <- ""
	for {
		select {
		case <-update:
			latestVer, err := Fetch.GitHubDownloadGet()
			if err != nil {
				DebugLogger.Println(err)
				continue
			}
			fmt.Println("latest version", latestVer)
			originVer := Fetch.GetVersion()
			fmt.Println("origin version", originVer)
			if latestVer != originVer {
				Fetch.Update()
			}
		case <-time.After(checkFreq):
			var now = time.Now()
			//if now.Day() != timeRecord.Day() {
			update <- ""
			//}
			//timeRecord = now
			DebugLogger.Println(now, "  time tick")
		}
	}
}
