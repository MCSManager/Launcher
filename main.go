package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/MCSManager/Launcher/lang"
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var webProcess *ProcessMgr
var daemonProcess *ProcessMgr

var infoView *tview.TextView
var operationView *tview.List
var outputWebView *tview.TextView
var outputDaemonView *tview.TextView

func main() {

	app := tview.NewApplication()

	lang.InitTranslations()
	lang.SetLanguage("zh-CN")
	initOperationCommand()

	// globalBox := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	infoView = tview.NewTextView().SetText("欢迎使用 MCSManager 面板")
	initUIBox(infoView, "运行状态")

	outputWebView = tview.NewTextView()
	initUIBox(outputWebView, "Web 日志")

	outputDaemonView = tview.NewTextView()
	initUIBox(outputDaemonView, "Daemon 日志")

	operationView = tview.NewList()
	operationView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		app.Sync()
		return event
	})
	initUIBox(operationView, "操作")

	for _, item := range globalOperationCommand {
		operationView.AddItem(item.MainText, item.SecondaryText, item.Shortcut, item.Exec)
	}

	flex2 := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(infoView, 0, 1, false).
		AddItem(operationView, 0, 1, true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flex2, 0, 1, true).
		AddItem(outputDaemonView, 0, 1, true).
		AddItem(outputWebView, 0, 1, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}

func onCommand(cmd string) {
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
