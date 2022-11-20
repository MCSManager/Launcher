package cmd

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

type ProcessMgr struct {
	Path     string
	Args     []string
	Started  bool
	startErr chan error
	exited   chan error
}

func NewProcessMgr(path string, args ...string) *ProcessMgr {
	return &ProcessMgr{Path: path, Args: args, startErr: make(chan error), exited: make(chan error)}
}

// ListenStop 监听程序停止运行
func (pm *ProcessMgr) ListenStop(callback func(err error)) {
	go func() {
		for {
			select {
			case err := <-pm.exited:
				callback(err)
				pm.Started = false
			}
		}
	}()
}

func (pm *ProcessMgr) Start() error {
	go pm.run()
	return <-pm.startErr
}

func (pm *ProcessMgr) run() {
	cmder := exec.Command(pm.Path, pm.Args...)
	err := cmder.Start()
	pm.startErr <- err
	pm.Started = true
	pm.exited <- cmder.Wait()
}

func (pm *ProcessMgr) End() error {
	processName := filepath.Base(pm.Path)

	var cmder *exec.Cmd
	if runtime.GOOS != "windows" {
		cmder = exec.Command("killall", "-9", processName)
	} else {
		//todo 待验证
		cmder = exec.Command("taskkill", "/F", "/im", processName)
	}
	err := cmder.Run()
	pm.Started = false
	return err
}
