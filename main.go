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

		if !scanner.Scan() {
			break
		}
		command := scanner.Text()
		onCommand(command)

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printPanelStatus() {
	fmt.Print(lang.T("PanelStatus"))
	if webProcess != nil && webProcess.Started {
		fmt.Println(color.GreenString(lang.T("running")))
	} else {
		fmt.Println(color.RedString(lang.T("stopped")))
	}
}

func helpInfo() {
	color.Green(lang.T("WelcomeTip"))

	fmt.Println()
	printPanelStatus()

	fmt.Println()
	fmt.Println(lang.T("HelpList"))

	fmt.Println()
	fmt.Println(color.HiYellowString(lang.T("PleaseInput")))

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
	if cmd == "4" {
		fmt.Println(color.HiGreenString(lang.T("AdvancedOptionHelp")))
		return
	}
	if cmd == "p1" {
		outputSubProcessLog(webProcess)
		return
	}
	if cmd == "p2" {
		outputSubProcessLog(daemonProcess)
		return
	}
	if cmd == "e" {
		stopPanel()
		os.Exit(0)
		return
	}
	fmt.Println(color.HiYellowString(lang.T("UnknownCommand")))
}

func stopPanel() {
	if webProcess != nil && daemonProcess != nil {
		webProcess.End()
		daemonProcess.End()
		fmt.Println(color.GreenString(lang.T("CommandSendSuccess")))
		return
	}
	fmt.Println(color.HiYellowString(lang.T("NotRunning")))
}

func startPanel() {

	if daemonProcess != nil || webProcess != nil {
		println(color.HiYellowString(lang.T("NotStopped")))
		return
	}

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

	printPanelStatus()

	wg.Wait()
	webProcess = nil
	daemonProcess = nil
}

func outputSubProcessLog(process *ProcessMgr) {
	if process == nil {
		return
	}
	process.IsOpenStdout = !process.IsOpenStdout
}

func startPanelProcess(cmd string, args ...string) *ProcessMgr {
	process := NewProcessMgr("/", cmd, "exit1", args...)
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
			fmt.Println(color.RedString("PROCESS ERR: "), out)
		}
	}()

	<-process.StartedEvent
	return process
}
