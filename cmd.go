package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type ProcessMgr struct {
	Path          string
	Args          []string
	Started       bool
	listenError   func(error)
	listenStdout  func(string)
	listenStarted func()
	listenExit    func()
	Cwd           string
	StartCount    int
	StopCommand   string
	IsOpenStdout  bool

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	cmder  *exec.Cmd
	wg     *sync.WaitGroup
}

func NewProcessMgr(workDir string, path string, stopCommand string, args ...string) *ProcessMgr {
	return &ProcessMgr{
		Path:        path,
		Args:        args,
		Cwd:         workDir,
		StopCommand: stopCommand,
		wg:          &sync.WaitGroup{},
	}
}

func (pm *ProcessMgr) ListenError(fn func(error)) {
	pm.listenError = fn
}
func (pm *ProcessMgr) ListenStdout(fn func(string)) {
	pm.listenStdout = fn
}
func (pm *ProcessMgr) ListenStarted(fn func()) {
	pm.listenStarted = fn
}

func (pm *ProcessMgr) ListenExit(fn func()) {
	pm.listenExit = fn
}

func (pm *ProcessMgr) Start() {
	go pm.run()
}

func (pm *ProcessMgr) run() error {
	os.Chdir(pm.Cwd)
	pm.StartCount += 1
	pm.cmder = exec.Command(pm.Path, pm.Args...)

	if runtime.GOOS == "windows" {
		pm.cmder.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	var err error
	pm.stdin, err = pm.cmder.StdinPipe()
	if pm.isError(err) {
		return err
	}
	pm.stdout, err = pm.cmder.StdoutPipe()
	if pm.isError(err) {
		return err
	}

	pm.stderr, err = pm.cmder.StderrPipe()
	if pm.isError(err) {
		return err
	}

	err = pm.cmder.Start()
	if pm.isError(err) {
		return err
	}

	pm.wg.Add(2)
	go pm.readStream(pm.stdout)
	go pm.readStream(pm.stderr)

	pm.Started = true
	pm.cmder.Wait()
	pm.Started = false
	if pm.listenExit != nil {
		pm.listenExit()
	}
	pm.wg.Wait()
	pm.close()
	return nil
}

func (pm *ProcessMgr) Write(text string) error {
	_, err := pm.stdin.Write([]byte(text + "\n"))
	go pm.ExitCheck()
	return err
}

func (pm *ProcessMgr) readStream(stream io.ReadCloser) {
	reader := bufio.NewReader(stream)
	for {
		buf := make([]byte, 512)
		n, err := reader.Read(buf)
		if err != nil || err == io.EOF {
			break
		}
		if pm.listenStdout != nil {
			pm.listenStdout(string(buf[:n]))
		}
	}
	defer stream.Close()
	defer pm.wg.Done()
}

func (pm *ProcessMgr) close() {
	// TODO
}

func (pm *ProcessMgr) isError(err error) bool {
	if err != nil {
		if pm.listenError != nil {
			pm.listenError(err)
		}
		pm.close()
		return true
	}
	return false
}

func (pm *ProcessMgr) End() error {
	if pm.cmder == nil || pm.stdin == nil {
		return nil
	}
	defer pm.stdin.Close()
	defer pm.stdout.Close()
	_, err := pm.stdin.Write([]byte(pm.StopCommand + "\n"))

	go pm.ExitCheck()
	return err
}

func (pm *ProcessMgr) ExitCheck() error {
	time.Sleep(5 * time.Second)
	tmpStartCount := pm.StartCount
	if pm.Started && pm.StartCount == tmpStartCount {
		pid := pm.cmder.Process.Pid
		// Only Windows support taskkill
		cmder := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
		return cmder.Run()
	}
	return nil
}
