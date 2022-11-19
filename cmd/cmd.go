package cmd

import (
	"os/exec"
	"runtime"
)

func Start(path string, errChan chan error, args ...string) {
	cmder := exec.Command(path, args...)
	err := cmder.Start()
	errChan <- err
}

func End(processName string, errChan chan error) {
	var cmder *exec.Cmd
	if runtime.GOOS != "windows" {
		cmder = exec.Command("killall", "-9", processName)
	} else {
		//todo 待验证
		cmder = exec.Command("taskkill", "/F", "/im", processName)
	}
	err := cmder.Start()
	errChan <- err
}
