package Fetch

import (
	. "../Debug"
	"../Envpath"
	"os/exec"
	"runtime"
)

var cmd *exec.Cmd = nil

const osType = runtime.GOOS

func kill() {
	if osType != "linux" {
		return
	}
	if cmd == nil {
		return
	}
	err := cmd.Process.Kill()
	if err != nil {
		DebugLogger.Println(err)
	}
	err = cmd.Process.Release()
	if err != nil {
		DebugLogger.Println(err)
	}
	cmd = nil
}

func start() {
	if osType != "linux" {
		return
	}
	if cmd != nil {
		DebugLogger.Println("YoC process is not killed normally ")
	}
	filepath := Envpath.GetAppDir() + "/YoC/YoC.bin"
	cmd = exec.Command(filepath, "")
	cmd.Stderr = GetLogWriter()
	cmd.Stdout = GetLogWriter()
	err := cmd.Start()
	if err != nil {
		DebugLogger.Println(err)
	}
	go func() {
		DebugLogger.Println("the process runs successfully")
		DebugLogger.Println("listening now . . .")
		err := cmd.Wait()
		DebugLogger.Println(err)
	}()
}

func Watch() {
	if cmd == nil {
		start()
		return
	} else {
		stat := cmd.ProcessState
		if stat == nil {
			return
		} else {
			DebugLogger.Println(stat.String())
			DebugLogger.Println(stat.UserTime())
			kill()
			start()
		}
	}
}
