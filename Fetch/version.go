package Fetch

import (
	"../Common/envpath"
	"bufio"
	"encoding/json"
	"io"
	"os"
)
import . "../Log"

func GetVersion() (ver string) {
	//Get info file path
	var dirPath = envpath.GetAppDir()
	dirPath, err := envpath.GetParentDir(dirPath)
	if err != nil {
		Log.Fatal(err)
	}
	dirPath, err = envpath.GetSubPath(dirPath, "YoC")
	if err != nil {
		Log.Fatal(err)
	}
	filePath, err := envpath.GetSubFile(dirPath, "YoC.info")
	if err != nil {
		Log.Println(err)
		return "0.0.0"
	}
	//Read file
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		Log.Fatal(err)
	}
	scanner := bufio.NewReader(file)
	bytes, err := scanner.ReadBytes('\n')
	if err != nil && err != io.EOF {
		Log.Fatal(err)
	}
	var global = make(map[string]string)
	err = json.Unmarshal(bytes, &global)
	if err != nil {
		Log.Println(err)
		return "0.0.0"
	}
	ver = global["Version"]
	return ver
}
