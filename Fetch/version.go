package Fetch

import (
	"../Envpath"
	"bufio"
	"encoding/json"
	"io"
	"os"
)
import . "../Debug"

func GetVersion() (ver string) {
	//Get info file path
	var dirPath = Envpath.GetAppDir()
	dirPath, err := Envpath.GetParentDir(dirPath)
	if err != nil {
		DebugLogger.Fatal(err)
	}
	dirPath, err = Envpath.GetSubPath(dirPath, "YoC")
	if err != nil {
		if os.IsNotExist(err) {
			DebugLogger.Println(err)
			err = Envpath.CheckMakeDir(dirPath)
			if err != nil {
				DebugLogger.Fatal(err)
			}
		} else {
			DebugLogger.Fatal(err)
		}
	}
	filePath, err := Envpath.GetSubFile(dirPath, "YoC.info")
	if err != nil {
		DebugLogger.Println(err)
		return "0.0.0"
	}
	//Read file
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		DebugLogger.Println(err)
		return "0.0.0"
	}
	scanner := bufio.NewReader(file)
	bytes, err := scanner.ReadBytes('\n')
	if err != nil && err != io.EOF {
		DebugLogger.Fatal(err)
	}
	var global = make(map[string]string)
	err = json.Unmarshal(bytes, &global)
	if err != nil {
		DebugLogger.Println(err)
		return "0.0.0"
	}
	if ver, ok := global["Version"]; ok {
		DebugLogger.Println("current version", ver)
		return ver
	} else {
		DebugLogger.Println("current version not exist")
		return "0.0.0"
	}
}
