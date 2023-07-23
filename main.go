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
	if cmd == "start" {
		go startPanel()
		return
	}

	if cmd == "exit" {
		stopPanel()
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

	webProcess = NewProcessMgr("/", "ping", "exit1", "www.baidu.com")
	webProcess.Start()
	daemonProcess = NewProcessMgr("/", "ping", "exit1", "www.baidu.com")
	daemonProcess.Start()

	var wg sync.WaitGroup
	wg.Add(2)

	daemonProcess.ListenExit(func() {
		wg.Done()
		webProcess.End()
	})

	webProcess.ListenExit(func() {
		wg.Done()
		daemonProcess.End()
	})

	// daemonProcess.ListenStdout(func(text string) {
	// 	fmt.Println("XZX: " + text)
	// })

	// webProcess.ListenStdout(func(text string) {
	// 	fmt.Println("AAA: " + text)
	// })

	wg.Wait()
	fmt.Println("程序退出")
	webProcess = nil
	daemonProcess = nil
}

func outputSubProcessLog(process *ProcessMgr) {
	if process == nil {
		return
	}
	process.IsOpenStdout = !process.IsOpenStdout
}
