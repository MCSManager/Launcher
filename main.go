package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"mcsmanager.com/desktop-app/cmd"
)

//go build -ldflags -H=windowsgui main.go

func main() {
	fontPath := "./config/msyh.ttc"
	os.Setenv("FYNE_FONT", fontPath)

	a := app.New()
	w := a.NewWindow("MCSManager 面板管理小程序")

	w.Resize(fyne.Size{Width: 280, Height: 360})

	// 数据源双向绑定
	statusLabelText := binding.NewString()
	statusLabelText.Set("未运行")

	statusLabel := widget.NewLabelWithData(statusLabelText)

	tipLabel := widget.NewLabel("请打开浏览器访问 http://localhost:23333/ 来使用。")

	// exitTipLabel := widget.NewLabel("必须先点击“关闭”按钮才可关闭窗口，否则可能会有数据损坏。")
	statusTipLabel := widget.NewLabel("状态")
	exitTipLabel := canvas.NewText("必须先点击“关闭”按钮才可关闭窗口，否则可能会有数据损坏。", color.Black)
	exitTipLabel.TextSize = 12

	//守护进程管理
	daemon := cmd.NewProcessMgr("ping", "www.baidu.com")

	btn := widget.NewButton("启动", nil)
	btnToggle := false

	//监听程序运行状态
	daemon.ListenStop(func(err error) {
		content := "已停止运行"
		if err != nil {
			content = fmt.Sprintf("%s\nerror: %s", content, err.Error())
		}
		statusLabelText.Set(content)
		btn.SetText("启动")
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
				statusLabelText.Set("正在运行")
			}
			btn.SetText("停止")
		} else { //停止
			if !daemon.Started {
				return
			}
			btn.SetText("停止中...")
			if err = daemon.End(); err != nil {
				btn.SetText(fmt.Sprintf("停止失败,error:%s", err.Error()))
			} else {
				statusLabelText.Set("未运行")
			}
		}
	}

	btnContainer := container.New(
		layout.NewMaxLayout(),
		btn,
	)

	firstLine := container.New(
		layout.NewHBoxLayout(),
		statusTipLabel,
		statusLabel,
	)

	content := container.New(layout.NewVBoxLayout(), firstLine, tipLabel, layout.NewSpacer(), exitTipLabel, btnContainer)

	w.SetContent(content)

	w.ShowAndRun()
}
