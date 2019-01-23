package main

import (
	"./Log"
	"log"
)

func main() {
	err := FetchLog.LogInit()
	if err != nil {
		log.Fatal(err)
	}
	FetchLog.LogExit(nil)
}