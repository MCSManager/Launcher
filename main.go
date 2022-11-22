package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"mcsmanager.com/desktop-app/cmd"
	"mcsmanager.com/desktop-app/uiw"
	"mcsmanager.com/desktop-app/utils"
)

//go build -ldflags -H=windowsgui main.go

func main() {

	STOPPED_TEXT := "状态: 未运行"
	STARTED_TEXT := "状态: 正在运行"

	fontPath := "./config/msyh.ttc"
	os.Setenv("FYNE_FONT", fontPath)

	app := app.New()
	window := app.NewWindow("MCSManager Launcher")

	window.Resize(fyne.Size{Width: 320, Height: 260})
	window.SetCloseIntercept(func() {
		fmt.Println("正在关闭窗口...")
		dialog.ShowConfirm("警告", "确定要退出程序吗？", func(b bool) {
			if b {
				os.Exit(0)
			}
		}, window)
	})

	statusLabel := uiw.NewMyLabel(STOPPED_TEXT)
	statusLabel.SetFontSize(12)
	tipLabel := uiw.NewMyLabel("端口: 23333")
	tipLabel.SetFontSize(12)
	tipLabelWrapper := container.New(layout.NewHBoxLayout(), tipLabel.Canvas)

	// exitTipLabel := uiw.NewMyLabel("必须点击关闭后台程序才可关闭窗口。")
	// exitTipLabel.SetFontSize(11)
	// exitTipLabel.SetColor(&color.RGBA{1, 2, 3, 200})

	//守护进程管理
	daemon := cmd.NewProcessMgr("./out/test")

	operationButton := widget.NewButton("启动后台程序", nil)
	// btnColor := canvas.NewRectangle(utils.BLUE)
	btnWrapper := container.New(
		layout.NewMaxLayout(),
		// btnColor,
		operationButton,
	)
	openBrower := widget.NewButton("访问面板", func() {
		fmt.Println("打开浏览器")
	})
	btnToggle := false

	//监听程序运行状态
	daemon.ListenStop(func(err error) {

		content := "已停止运行"
		if err != nil {
			content = fmt.Sprintf("%s\nerror: %s", content, err.Error())
			fmt.Println(content)
		}
		operationButton.SetText("启动后台程序")
		btnToggle = false
	})

	operationButton.OnTapped = func() {
		btnToggle = !btnToggle
		operationButton.Disable()
		defer operationButton.Enable()
		var err error
		if btnToggle { //启动
			if daemon.Started {
				return
			}
			operationButton.SetText("启动中...")
			if err = daemon.Start(); err != nil {
				operationButton.SetText(fmt.Sprintf("启动失败,error:%s", err.Error()))
			} else {
				statusLabel.SetText(STARTED_TEXT)
				statusLabel.SetColor(utils.GREEN)
				statusLabel.Canvas.Refresh()
			}
			operationButton.SetText("停止后台程序")
		} else { //停止
			if !daemon.Started {
				return
			}
			operationButton.SetText("停止中...")
			if err = daemon.End(); err != nil {
				operationButton.SetText(fmt.Sprintf("停止失败,error:%s", err.Error()))
			} else {
				statusLabel.SetText(STOPPED_TEXT)
				statusLabel.SetColor(color.Black)
			}
		}
	}

	infoLabel := uiw.NewMyLabel("MCSManager 面板启动器")
	infoLabel.SetFontSize(12)

	paddingContainer1 := container.New(layout.NewPaddedLayout(), infoLabel.Canvas)
	paddingContainer2 := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(), statusLabel.Canvas, tipLabelWrapper))
	paddingContainer3 := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(), btnWrapper, openBrower))

	content := container.New(layout.NewVBoxLayout(), paddingContainer1, layout.NewSpacer(), paddingContainer2, paddingContainer3)

	window.SetContent(container.New(layout.NewPaddedLayout(), content))

	window.ShowAndRun()
}
