package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"mcsmanager.com/desktop-app/cmd"
	"mcsmanager.com/desktop-app/uiw"
)

//go build -ldflags -H=windowsgui main.go

func main() {

	STOPPED_TEXT := "未运行"
	STARTED_TEXT := "运行中"

	fontPath := "./config/msyh.ttc"
	os.Setenv("FYNE_FONT", fontPath)

	a := app.New()
	w := a.NewWindow("MCSManager 面板管理小程序")

	w.Resize(fyne.Size{Width: 280, Height: 360})

	statusLabel := uiw.NewMyLabel(STOPPED_TEXT)

	tipLabel := uiw.NewMyLabel("请打开浏览器访问 http://localhost:23333/ 来使用。")

	// exitTipLabel := uiw.NewMyLabel("必须先点击“关闭”按钮才可关闭窗口，否则可能会有数据损坏。")
	statusTipLabel := uiw.NewMyLabel("状态")

	exitTipLabel := canvas.NewText("必须点击关闭后台程序才可关闭窗口，否则可能会有数据损坏。", &color.RGBA{1, 2, 3, 200})
	exitTipLabel.TextSize = 11

	//守护进程管理
	daemon := cmd.NewProcessMgr("ping", "www.baidu.com")

	btn := widget.NewButton("启动后台程序", nil)
	btnToggle := false

	//监听程序运行状态
	daemon.ListenStop(func(err error) {
		content := "已停止运行"
		if err != nil {
			content = fmt.Sprintf("%s\nerror: %s", content, err.Error())
			fmt.Println(content)
		}
		btn.SetText("启动后台程序")
		btnToggle = false
	})

	btn.OnTapped = func() {
		btnToggle = !btnToggle
		btn.Disable()
		defer btn.Enable()
		var err error
		if btnToggle { //启动
			if daemon.Started {
				return
			}
			btn.SetText("启动中...")
			if err = daemon.Start(); err != nil {
				btn.SetText(fmt.Sprintf("启动失败,error:%s", err.Error()))
			} else {
				statusLabel.SetText(STARTED_TEXT)
				statusLabel.SetColor(color.RGBA{0, 244, 0, 255})
				// statusLabel.Canvas.Resize(fyne.Size{Height: 300, Width: 100})
				statusLabel.Canvas.Refresh()
			}
			btn.SetText("停止后台程序")
		} else { //停止
			if !daemon.Started {
				return
			}
			btn.SetText("停止中...")
			if err = daemon.End(); err != nil {
				btn.SetText(fmt.Sprintf("停止失败,error:%s", err.Error()))
			} else {
				statusLabel.SetText(STOPPED_TEXT)
				statusLabel.SetColor(color.Black)
			}
		}
	}

	btnContainer := container.New(
		layout.NewMaxLayout(),
		btn,
	)

	firstLine := container.New(
		layout.NewHBoxLayout(),
		statusTipLabel.Canvas,
		statusLabel.Canvas,
	)

	content := container.New(layout.NewVBoxLayout(), firstLine, tipLabel.Canvas, layout.NewSpacer(), exitTipLabel, btnContainer)

	w.SetContent(content)

	w.ShowAndRun()
}
