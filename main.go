package main

import (
	. "./Debug"
	"./Fetch"
	"fmt"
	"log"
	"time"
)

const watchFreq = time.Minute
const checkFreq = time.Minute * 5
const tickFreq = time.Hour

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
		case <-time.After(watchFreq):
			Fetch.Watch()
		case <-time.After(checkFreq):
			update <- ""
		case <-time.After(tickFreq):
			DebugLogger.Println(time.Now(), "  time tick")
		}
	}
}
