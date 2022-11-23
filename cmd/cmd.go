package cmd

import (
	"io"
	"os/exec"
)

type ProcessMgr struct {
	Path     string
	Args     []string
	Started  bool
	stdin    io.WriteCloser
	startErr chan error
	exited   chan error
	cmder    *exec.Cmd
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
				pm.Started = false
				callback(err)
			}
		}
	}()
}

func (pm *ProcessMgr) Start() error {
	go pm.run()
	return <-pm.startErr
}

func (pm *ProcessMgr) run() {
	pm.cmder = exec.Command(pm.Path, pm.Args...)
	var err error
	pm.stdin, err = pm.cmder.StdinPipe()
	if err != nil {
		pm.startErr <- err
		return
	}
	err = pm.cmder.Start()
	pm.startErr <- err
	pm.Started = true
	pm.exited <- pm.cmder.Wait()
}

func (pm *ProcessMgr) End() error {
	if pm.cmder == nil || pm.stdin == nil {
		return nil
	}
	defer pm.stdin.Close()

	_, err := pm.stdin.Write([]byte("exit\n"))
	return err
}
