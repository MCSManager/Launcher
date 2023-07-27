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

	// go func() {
	// 	for {
	clearTerminal()

	go startPanel()
	time.Sleep(1000 * time.Millisecond)

	// totalSecond = totalSecond + 1
	// days, hours, minutes, remainingSeconds := formatDuration(totalSecond)

	fmt.Println(color.HiGreenString("---------------------------"))
	fmt.Println(color.HiGreenString(lang.T("WelcomeTip")))
	fmt.Println(color.CyanString(lang.T("SoftwareInfo")))
	fmt.Println(color.HiGreenString("---------------------------"))
	fmt.Println()
	fmt.Println(color.HiGreenString(lang.T("PanelStatus")) + getPanelStatusText())
	// fmt.Println(color.HiGreenString(lang.FT("RunTime", map[string]interface{}{
	// 	"Time": color.HiYellowString(lang.FT("TimeText", map[string]interface{}{
	// 		"D": days,
	// 		"H": hours,
	// 		"M": minutes,
	// 		"S": remainingSeconds,
	// 	})),
	// })))
	fmt.Println()

	fmt.Println(lang.FT("Address", map[string]interface{}{
		"Url": color.HiYellowString(defaultHttpAddr),
	}))

	fmt.Println(color.WhiteString(lang.T("ExitTip")))
	fmt.Println()

	// 		time.Sleep(1000 * time.Millisecond)
	// 	}
	// }()

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

func clearTerminal() {
	// c := exec.Command("clear")
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func onCommand(cmd string) {
	if cmd == "s" {
		go startPanel()
		return
	}

	if cmd == "c" {
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
		os.Exit(0)
	}
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

	webProcess = nil
	daemonProcess = nil
}
