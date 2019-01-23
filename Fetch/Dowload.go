package Fetch

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
import . "../Log"

func DowloadFile(filePath string, fileUrl string) error {
	var client= http.DefaultClient
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

func GitHubDownloadGet(fileUrl string) (string, error) {
	var client = http.DefaultClient
	client.Timeout = time.Second * 60
	resp, err := client.Get(fileUrl)
	if err != nil {
		Log.Println(err)
		return "", err
	}
	return resp.Header.Get("Location"), nil
}
