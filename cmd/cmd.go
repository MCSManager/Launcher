package cmd

import (
	"os/exec"
)

type ProcessMgr struct {
	Path     string
	Args     []string
	Started  bool
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
	pm.cmder = exec.Command(pm.Path, pm.Args...)
	err := pm.cmder.Start()
	pm.startErr <- err
	pm.Started = true
	pm.exited <- pm.cmder.Wait()
}

func (pm *ProcessMgr) End() error {
	if pm.cmder == nil {
		return nil
	}
	//err := pm.cmder.Process.Signal(syscall.SIGTERM)
	stdin, err := pm.cmder.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	_, err = stdin.Write([]byte("exit\n"))
	return err
}
