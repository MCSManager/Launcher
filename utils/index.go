package utils

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
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

// Open calls the OS default program for uri
func Open(uri string) error {
	//run, ok := commands[runtime.GOOS]
	//if !ok {
	//return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	//}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd.exe", "/c", "start", uri)
		return cmd.Start()
	}
	return errors.New("not support")
}
