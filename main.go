package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/MCSManager/Launcher/lang"
	"github.com/fatih/color"
)

var webProcess *ProcessMgr
var daemonProcess *ProcessMgr

func main() {

	lang.InitTranslations()
	lang.SetLanguage("zh-CN")
	helpInfo()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		command := scanner.Text()
		onCommand(command)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func helpInfo() {
	color.Green(lang.T("WelcomeTip"))

	fmt.Println()
	fmt.Print(lang.T("PanelStatus"))
	fmt.Println(color.GreenString(lang.T("running")))

	fmt.Println()
	fmt.Println(lang.T("HelpList"))

	fmt.Println()
	fmt.Println(color.YellowString(lang.T("PleaseInput")))

}

func onCommand(cmd string) {
	if cmd == "h" {
		helpInfo()
		return
	}
	if cmd == "1" {
		println("成功！")
		return
	}
	if cmd == "2" {
		go startPanel()
		return
	}
	if cmd == "3" {
		go stopPanel()
		return
	}

}

func stopPanel() {
	if webProcess != nil && daemonProcess != nil {
		webProcess.End()
		daemonProcess.End()
		return
	}
	fmt.Println("The Panel is not running")
}

func startPanel() {

	webProcess = startPanelProcess("ping", "www.baidu.com")
	daemonProcess = startPanelProcess("ping", "www.google.com")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		<-daemonProcess.ExitEvent
		webProcess.End()
		wg.Done()
	}()
	go func() {
		<-webProcess.ExitEvent
		daemonProcess.End()
		wg.Done()
	}()

	wg.Wait()
	webProcess = nil
	daemonProcess = nil
}

func startPanelProcess(cmd string, args ...string) *ProcessMgr {
	process := NewProcessMgr("/", cmd, "exit", args...)
	process.Start()

	go func() {
		for {
			out, ok := <-process.StdoutEvent
			if !ok {
				break
			}
			fmt.Print(out)
		}
	}()

	go func() {
		for {
			out, ok := <-process.ErrEvent
			if !ok {
				break
			}
			fmt.Println("错误: ", out)
		}
	}()

	ok := <-process.StartedEvent
	fmt.Println("启进程结果如下：", ok)

	return process
}
