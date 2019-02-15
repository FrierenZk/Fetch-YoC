package Fetch

import (
	"../Envpath"
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
import . "../Debug"

const githubPage = "https://github.com/"
const projectName = "FrierenZk/YoC"
const projectReleases = projectName + "/releases"
const githubProjectPage = githubPage + projectReleases

var download, filename string

func Update() {
	filepath := Envpath.GetAppDir() + "/tmp"
	err := Envpath.CheckMakeDir(filepath)
	if err != nil {
		DebugLogger.Fatal(err)
	}
	err = downloadFile(filepath)
	if err != nil {
		DebugLogger.Println(err)
		return
	}
	//TODO check and kill the YoC process
	copyFile(filepath + "/" + filename)
}

func GitHubDownloadGet() (ver string, err error) {
	fileUrl := githubProjectPage + "/latest"
	var client, downloadLink = http.DefaultClient, ""
	client.Timeout = time.Second * 60
	resp, err := client.Get(fileUrl)
	if err != nil {
		DebugLogger.Println(err)
		return "", err
	}
	scanner := bufio.NewReader(resp.Body)
	bytes, err := scanner.ReadBytes('\n')
	for err != io.EOF || err == nil {
		line := string(bytes)
		if strings.Contains(line, projectReleases+"/download") {
			downloadLink = line
			//fmt.Println(line)
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

func getVersion(str string) (ver string) {
	strArr := strings.Split(str, "/")
	var i int
	for i = 0; i < len(strArr); i++ {
		if strArr[i] == "download" {
			break
		}
	}
	if len(strArr) < i+2 {
		DebugLogger.Println("link error :", str)
		return ""
	}
	ver = strArr[i+1]
	filename = strArr[i+2]
	DebugLogger.Println("find latest version : ", ver)
	return ver
}

func downloadFile(filePath string) error {
	var fileUrl = download
	var client = http.DefaultClient
	client.Timeout = time.Second * 60
	resp, err := client.Get(fileUrl)
	if err != nil {
		DebugLogger.Println(err)
		return err
	}
	disposition := resp.Header.Get("Content-Disposition")
	filename := strings.Split(disposition, "filename=")
	if len(filename) > 1 {
		filePath += "/" + filename[1]
	} else {
		err = errors.New("can't get filename")
		DebugLogger.Println(err)
		return err
	}
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		DebugLogger.Println(err)
		return err
	}
	scanner, writer, process := bufio.NewReader(resp.Body), bufio.NewWriter(file), int64(0)
	length, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	DebugLogger.Println("download file size", length)
	var readCount int64 = 0
	for readCount < length {
		data, err := scanner.ReadBytes(0)
		if err != nil && err != io.EOF {
			DebugLogger.Println(err)
			return err
		}
		nBytes := len(data)
		if nBytes > 0 {
			_, _ = writer.Write(data)
			readCount += int64(nBytes)
			currentProcess := readCount * 1000
			currentProcess /= length
			if currentProcess != process {
				fmt.Println(float64(process) / 10)
				process = currentProcess
			}
		}
	}
	_ = writer.Flush()

	if readCount >= length {
		DebugLogger.Println("download complete")
	} else {
		err = errors.New("download size error")
		DebugLogger.Println(err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
		_ = file.Close()
	}()
	return nil
}

func copyFile(src string) {
	var dst = Envpath.GetAppDir() + "/YoC"
	err := Envpath.CheckMakeDir(dst)
	if err != nil {
		DebugLogger.Println(err)
		return
	}
	dst += "/YoC.bin"
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	defer func() { _ = dstFile.Close() }()
	if err != nil {
		DebugLogger.Println(err)
		return
	}
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	defer func() { _ = srcFile.Close() }()
	if err != nil {
		DebugLogger.Println(err)
		return
	}
	n, err := io.Copy(dstFile, srcFile)
	DebugLogger.Println(n, "bytes copied")
	if err != nil {
		DebugLogger.Println(err)
	}
}
