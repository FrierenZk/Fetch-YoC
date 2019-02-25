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
	var update, watch, tick, start = time.NewTicker(updateFreq), time.NewTicker(watchFreq), time.NewTicker(tickFreq), time.NewTimer(time.Nanosecond)
	defer update.Stop()
	defer watch.Stop()
	defer tick.Stop()
	for {
		select {
		case <-start.C:
			//Work only once
			Update()
		case <-update.C:
			Update()
		case <-watch.C:
			Fetch.Watch()
		case <-tick.C:
			DebugLogger.Println(time.Now(), "  time tick")

		}
	}
}

func Update() {
	latestVer, err := Fetch.GitHubDownloadGet()
	if err != nil {
		DebugLogger.Println(err)
		return
	}
	fmt.Println("latest version", latestVer)
	originVer := Fetch.GetVersion()
	fmt.Println("origin version", originVer)
	if latestVer != originVer {
		Fetch.Update()
	}
}
