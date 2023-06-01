package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
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
	}
}

func (pm *ProcessMgr) Start() {
	go pm.run()
}

func (pm *ProcessMgr) run() {
	os.Chdir(pm.Cwd)
	fmt.Printf("Change CWD: %s %s\n", pm.Cwd, pm.Path)
	pm.StartCount += 1
	pm.cmder = exec.Command(pm.Path, pm.Args...)
	// if runtime.GOOS == "windows" {
	// 	pm.cmder.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// }
	var err error
	pm.stdin, err = pm.cmder.StdinPipe()
	if err != nil {
		pm.ErrEvent <- err
	}
	pm.stdout, err = pm.cmder.StdoutPipe()
	if err != nil {
		pm.ErrEvent <- err
	}

	pm.stderr, err = pm.cmder.StderrPipe()
	if err != nil {
		pm.ErrEvent <- err
	}

	err = pm.cmder.Start()
	println("start", err)
	if err != nil {
		pm.ErrEvent <- err
	}

	// Stdout
	go func() {
		defer pm.stdout.Close()
		reader := bufio.NewReader(pm.stdout)
		for {
			buf := make([]byte, 512)
			n, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					pm.IoErrEvent <- err
				}
			}
			pm.StdoutEvent <- string(string(buf[:n]))
		}
	}()

	// Stderr
	go func() {
		defer pm.stderr.Close()
		reader := bufio.NewReader(pm.stderr)
		for {
			buf := make([]byte, 512)
			n, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					pm.IoErrEvent <- err
				}
			}
			pm.StdoutEvent <- string(buf[:n])
		}
	}()

	// Stdin
	go func() {
		defer pm.stdin.Close()
		for {
			input := <-pm.StdinEvent
			io.WriteString(pm.stdin, input)
		}
	}()

	pm.Started = true
	println("启动成功")
	pm.StartedEvent <- nil
	pm.ExitEvent <- pm.cmder.Wait()
	pm.Started = false
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
		time.Sleep(6 * time.Second)
		fmt.Printf("Program kill, Status: %v", pm.Started)
		if pm.Started && pm.StartCount == tmpStartCount {
			pid := pm.cmder.Process.Pid
			fmt.Printf("Kill Program: taskkill /PID %d /T /F\n", pid)
			cmder := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
			// cmder.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			err := cmder.Run()
			if err != nil {
			}
		}
	}()
}
