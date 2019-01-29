package main

import (
	"./Fetch"
	"./Log"
	"fmt"
	"log"
)

func main() {
	err := FetchLog.LogInit()
	if err != nil {
		log.Fatal(err)
	}
	//fileUrl,err := Fetch.GitHubDownloadGet("https://github.com/shadowsocks/shadowsocks-android/releases/download/v4.7.0/shadowsocks--universal-4.7.0.apk")
	//err =Fetch.DowloadFile("E:\\GitHub\\","https://github.com/shadowsocks/shadowsocks-android/releases/download/v4.7.0/shadowsocks--universal-4.7.0.apk")
	ver,err := Fetch.GitHubDownloadGet()
	fmt.Println(ver)
	FetchLog.LogExit(nil)
}