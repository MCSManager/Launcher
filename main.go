package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/MCSManager/Launcher/lang"
	"github.com/fatih/color"
)

var webProcess *ProcessMgr
var daemonProcess *ProcessMgr
var totalSecond int64
var defaultHttpAddr = "http://127.0.0.1:23333"

func main() {

	lang.InitTranslations()
	lang.SetLanguage("zh-CN")

	clearTerminal()

	fmt.Println(color.HiGreenString("---------------------------"))
	fmt.Println(color.HiGreenString(lang.T("WelcomeTip")))
	fmt.Println(color.CyanString(lang.T("SoftwareInfo")))
	fmt.Println(color.HiGreenString("---------------------------"))

	go startPanel()

	// time.Sleep(5000 * time.Millisecond)

	// if !getPanelStatus() {
	// 	return
	// }

	fmt.Println()

	fmt.Println(lang.FT("Address", map[string]interface{}{
		"Url": color.HiYellowString(defaultHttpAddr),
	}))

	fmt.Println(color.WhiteString(lang.T("ExitTip")))
	fmt.Println()

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

func getPanelStatusText() string {
	var panelStatus = color.HiRedString(lang.T("stopped"))
	if webProcess != nil && daemonProcess != nil && webProcess.Started && daemonProcess.Started {
		panelStatus = color.GreenString(lang.T("running"))
	}
	return panelStatus
}

func getPanelStatus() bool {
	return webProcess != nil && daemonProcess != nil && webProcess.Started && daemonProcess.Started
}

func clearTerminal() {
	// c := exec.Command("clear")
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func onCommand(cmd string) {
	// if cmd == "start" {
	// 	go startPanel()
	// 	return
	// }

	if cmd == "stop" {
		stopPanel()
		return
	}

	logErr(lang.T("UnknownCommand") + cmd)
}

func stopPanel() {
	if webProcess != nil && daemonProcess != nil {
		logInfo(lang.T("stoppingPanel"))
		daemonProcess.End()
		webProcess.End()

	}
}

func startPanel() {
	if daemonProcess != nil || webProcess != nil {
		logErr(color.HiYellowString(lang.T("NotStopped")))
		return
	}

	daemonProcess = NewProcessMgr("/", "ping", "exit1", "www.google.com")
	webProcess = NewProcessMgr("/", "ping", "exit1", "www.google.com")

	var wg sync.WaitGroup
	wg.Add(2)

	daemonProcess.ListenExit(func() {
		go webProcess.End()
		logInfo(fmt.Sprintf("%s %s", lang.T("Daemon"), lang.T("DaemonProcessExit")))
		wg.Done()
	})

	webProcess.ListenExit(func() {
		go daemonProcess.End()
		logInfo(fmt.Sprintf("%s %s", lang.T("Web"), lang.T("WebProcessExit")))
		wg.Done()
	})

	daemonProcess.ListenError(func(err error) {
		logErr(fmt.Sprintf("%s %s", lang.T("Daemon"), err.Error()))
	})

	webProcess.ListenError(func(err error) {
		logErr(fmt.Sprintf("%s %s", lang.T("Web"), err.Error()))
	})

	// daemonProcess.ListenStdout(func(text string) {
	// 	logInfo(fmt.Sprintf("%s %s", lang.T("Daemon"), text))
	// })

	// webProcess.ListenStdout(func(text string) {
	// 	logInfo(fmt.Sprintf("%s %s", lang.T("Web"), text))
	// })

	err1 := daemonProcess.Start()
	if err1 != nil {
		onPanelExitEvent()
		return
	}

	time.Sleep(2 * time.Second)

	err2 := webProcess.Start()
	if err2 != nil {
		daemonProcess.End()
		onPanelExitEvent()
		return
	}

	wg.Wait()

	onPanelExitEvent()
}

func onPanelExitEvent() {
	// webProcess = nil
	// daemonProcess = nil

	fmt.Println()
	logErr(lang.T("ExitedTip"))
}
