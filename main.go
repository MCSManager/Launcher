package main

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"os"

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
	//fmt.Println("U %v", utils.IsFileExists(fontPath))
	a := app.New()
	w := a.NewWindow("MCSManager 面板管理小工具")

	w.Resize(fyne.Size{Width: 280, Height: 360})

	// 数据源双向绑定
	statusLabelText := binding.NewString()
	statusLabelText.Set("Initial value")

	statusLabel := widget.NewLabelWithData(statusLabelText)

	//守护进程管理
	daemon := cmd.NewProcessMgr("ping", "-c5", "www.baidu.com")

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
			btn.SetText("启动")
		}
	}

	// btn_color := canvas.NewRectangle(
	// 	color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	container1 := container.New(
		// layout of container
		layout.NewMaxLayout(),
		// first use btn color
		// btn_color,
		// 2nd btn widget
		btn,
	)

	content := container.New(layout.NewVBoxLayout(), statusLabel, layout.NewSpacer(), container1)

	w.SetContent(content)

	w.ShowAndRun()
}
