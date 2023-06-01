package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MCSManager/Launcher/lang"
	"github.com/fatih/color"
)

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

}

func startWebProcess() {
	startPanelProcess("bash")
}

// func startDaemonProcess(cmd string) {
// 	startPanelProcess("bash")
// }

func startPanelProcess(cmd string) {
	process := NewProcessMgr("/", cmd, "exit")
	process.Start()

	go func() {
		for {
			out := <-process.StdoutEvent
			fmt.Print("stdout: ", out)
		}
	}()

	go func() {
		for {
			out := <-process.IoErrEvent
			fmt.Print("IoErrEvent: ", out)
		}
	}()

	<-process.StartedEvent
	process.StdinEvent <- "ping www.baidu.com\n"

	err := <-process.ErrEvent
	fmt.Println("有错误：", err)
}
