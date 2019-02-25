package main

import (
	. "./Debug"
	"./Fetch"
	"fmt"
	"log"
	"time"
)

const watchFreq = time.Minute
const updateFreq = time.Minute * 5
const tickFreq = time.Hour

func main() {
	err := LogInit()
	if err != nil {
		log.Fatal(err)
	}
	defer LogExit(nil)
	var update, watch, tick = time.After(updateFreq), time.After(watchFreq), time.After(tickFreq)
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
		case <-watch:
			Fetch.Watch()
		case <-tick:
			DebugLogger.Println(time.Now(), "  time tick")
		}
	}
}
