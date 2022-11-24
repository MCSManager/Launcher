package utils

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

var commands = map[string]string{
	"windows": "cmd.exe /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var ERR_LOG_PATH = "launcher_err.log"

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func WriteErrLog(err string) {
	os.WriteFile(ERR_LOG_PATH, []byte(err), 0744)
}

func Open(uri string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd.exe", "/c", "start", uri)
		if runtime.GOOS == "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}
		return cmd.Start()
	}
	return errors.New("not support")
}
