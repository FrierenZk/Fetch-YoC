package Fetch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
import . "../Log"

const githubPage = "https://github.com/"
const projectName = "shadowsocks/shadowsocks-windows"
const projectReleases = projectName +"/releases"
const githubProjectPage = githubPage + projectReleases
var download string

func DowloadFile(filePath string) error {
	var fileUrl = download
	var client = http.DefaultClient
	client.Timeout = time.Second * 60
	//var reader io.Reader
	resp, err := client.Get(fileUrl)
	if err != nil {
		Log.Println(err)
		return err
	}
	disposition := resp.Header.Get("Content-Disposition")
	filename := strings.Split(disposition, "filename=")
	if len(filename) > 1 {
		filePath += filename[1]
	} else {
		err = errors.New("can't get filename")
		Log.Println(err)
		filePath += "YoC"
	}
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		Log.Println(err)
		return err
	}
	scanner, writer := bufio.NewReader(resp.Body), bufio.NewWriter(file)
	length, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	var readCount int64 = 0
	for readCount < length {
		data, err := scanner.Peek(4096)
		if err != nil && err != io.EOF {
			Log.Println(err)
			return err
		}
		_, _ = writer.Write(data)
		readCount += 4096
	}
	if readCount >= length {
		Log.Println("download complete")
	} else {
		err = errors.New("download size error")
		Log.Println(err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
		_ = file.Close()
	}()
	return nil
}

func GitHubDownloadGet() (ver string,err error) {
	fileUrl := githubProjectPage + "/latest"
	var client, downloadLink= http.DefaultClient, ""
	client.Timeout = time.Second * 60
	resp, err := client.Get(fileUrl)
	if err != nil {
		Log.Println(err)
		return "", err
	}
	scanner := bufio.NewReader(resp.Body)
	bytes, err := scanner.ReadBytes('\n')
	for err != io.EOF || err == nil {
		line := string(bytes)
		if strings.Contains(line, projectReleases+"/download") {
			downloadLink = line
			fmt.Println(line)
			break
		}
		bytes, err = scanner.ReadBytes('\n')
	}
	strArr := strings.Split(downloadLink, "\"")
	for _, str := range strArr {
		if strings.Contains(str, projectReleases) {
			download = str
			break
		}
	}
	ver = getVersion(download)
	download = githubPage + download
	return ver, nil
}

func getVersion(str string)(ver string) {
	strArr := strings.Split(str, "/")
	var i int
	for i = 0; i < len(strArr); i++ {
		if strArr[i] == "download" {
			break
		}
	}
	ver = strArr[i+1]
	Log.Println("find latest version : ", ver)
	return ver
}
