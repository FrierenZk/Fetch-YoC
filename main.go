package main

import (
	"./Debug"
	"./Fetch"
	"fmt"
	"log"
)

func main() {
	err := FetchLog.LogInit()
	if err != nil {
		log.Fatal(err)
	}
	//latestVer,err := Fetch.GitHubDownloadGet()
	//fmt.Println(ver)
	//err =Fetch.DowloadFile("E:/")
	originVer := Fetch.GetVersion()
	fmt.Println(originVer)
	FetchLog.LogExit(nil)
}
