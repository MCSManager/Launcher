package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type ProcessMgr struct {
	Path         string
	Args         []string
	Started      bool
	ExitEvent    chan error
	ErrEvent     chan error
	IoErrEvent   chan error
	StdoutEvent  chan string
	StdinEvent   chan string
	StartedEvent chan error
	Cwd          string
	StartCount   int
	StopCommand  string

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	cmder  *exec.Cmd
	wg     *sync.WaitGroup
}

func NewProcessMgr(workDir string, path string, stopCommand string, args ...string) *ProcessMgr {
	return &ProcessMgr{
		Path:         path,
		Args:         args,
		Cwd:          workDir,
		StdoutEvent:  make(chan string),
		StdinEvent:   make(chan string),
		IoErrEvent:   make(chan error),
		ExitEvent:    make(chan error),
		ErrEvent:     make(chan error),
		StartedEvent: make(chan error),
		StopCommand:  stopCommand,
		wg:           &sync.WaitGroup{},
	}
}

func (pm *ProcessMgr) Start() {
	go pm.run()
}

func (pm *ProcessMgr) run() {
	os.Chdir(pm.Cwd)
	fmt.Printf("启动程序: %s %s %v\n", pm.Cwd, pm.Path, pm.Args)
	pm.StartCount += 1
	pm.cmder = exec.Command(pm.Path, pm.Args...)
	// if runtime.GOOS == "windows" {
	// 	pm.cmder.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// }
	var err error
	pm.stdin, err = pm.cmder.StdinPipe()
	if err != nil {
		pm.ErrEvent <- err
		pm.close()
		return
	}
	pm.stdout, err = pm.cmder.StdoutPipe()
	if err != nil {
		pm.ErrEvent <- err
		pm.close()
		return
	}

	pm.stderr, err = pm.cmder.StderrPipe()
	if err != nil {
		pm.ErrEvent <- err
		pm.close()
		return
	}

	err = pm.cmder.Start()

	if err != nil {
		pm.ErrEvent <- err
		pm.close()
		return
	}

	pm.StartedEvent <- nil

	// Stdout and Stderr
	pm.wg.Add(2)
	go pm.readStream(pm.stdout)
	go pm.readStream(pm.stderr)

	// Stdin
	go func() {
		defer pm.stdin.Close()
		for {
			input, ok := <-pm.StdinEvent
			if !ok {
				break
			}
			io.WriteString(pm.stdin, input)
		}
	}()

	pm.Started = true
	println("进程完全启动成功")
	pm.ExitEvent <- pm.cmder.Wait()
	pm.Started = false
	pm.wg.Wait()
	pm.close()
	println("进程完全退出成功")
}

func (pm *ProcessMgr) readStream(stream io.ReadCloser) {

	reader := bufio.NewReader(stream)
	for {
		buf := make([]byte, 512)
		n, err := reader.Read(buf)
		if err != nil || err == io.EOF {
			break
		}
		pm.StdoutEvent <- string(buf[:n])
	}
	defer stream.Close()
	defer pm.wg.Done()
}

func (pm *ProcessMgr) close() {
	close(pm.StartedEvent)
	close(pm.ExitEvent)
	close(pm.StdinEvent)
	close(pm.ErrEvent)
	close(pm.IoErrEvent)
	close(pm.StdoutEvent)
}

func (pm *ProcessMgr) End() error {
	if pm.cmder == nil || pm.stdin == nil {
		return nil
	}
	defer pm.stdin.Close()
	defer pm.stdout.Close()
	_, err := pm.stdin.Write([]byte(pm.StopCommand + "\n"))
	pm.ExitCheck()
	return err
}

func (pm *ProcessMgr) ExitCheck() {
	go func() {
		tmpStartCount := pm.StartCount
		time.Sleep(1 * time.Second)
		if pm.Started && pm.StartCount == tmpStartCount {
			pid := pm.cmder.Process.Pid
			// Only Windows support taskkill
			cmder := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
			cmder.Run()
		}
	}()
}
