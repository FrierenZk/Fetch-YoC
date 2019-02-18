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
	err = cmd.Wait()
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
	var err error
	cmd.Stderr = GetLogWriter()
	cmd.Stdout = GetLogWriter()
	if err != nil {
		DebugLogger.Println(err)
	}
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Start()
	if err != nil {
		DebugLogger.Println(err)
	}
}
