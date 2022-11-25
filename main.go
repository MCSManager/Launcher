package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/MCSManager/Launcher/cmd"
	"github.com/MCSManager/Launcher/uiw"
	"github.com/MCSManager/Launcher/utils"
)

type WebConfig struct {
	HttpPort int `json:"httpPort"`
}

func main() {

	var webConfig WebConfig

	STOPPED_TEXT := "状态: 未运行"
	STARTED_TEXT := "状态: 正在运行"

	if utils.IsFileExists("C:/Windows/Fonts/msyh.ttc") {
		os.Setenv("FYNE_FONT", "C:/Windows/Fonts/msyh.ttc")
	} else {
		os.Setenv("FYNE_FONT", "./config/msyh.ttc")
	}

	app := app.New()
	window := app.NewWindow("MCSManager Launcher")
	window.Resize(fyne.Size{Width: 320, Height: 260})
	window.SetFullScreen(false)
	window.SetFixedSize(true)

	infoLabel := uiw.NewMyLabel("MCSManager 管理面板启动器")
	infoLabel.SetFontSize(12)

	statusLabel := uiw.NewMyLabel(STOPPED_TEXT)
	statusLabel.SetFontSize(12)

	tipLabel := uiw.NewMyLabel("")
	tipLabel.SetFontSize(12)
	tipLabelWrapper := container.New(layout.NewHBoxLayout(), tipLabel.Canvas)
	operationButton := widget.NewButton("启动后台程序", nil)
	btnWrapper := container.New(
		layout.NewMaxLayout(),
		operationButton,
	)

	WEB_CONFIG_FILE_PATH := "./mcsmanager/web/data/SystemConfig/config.json"
	if utils.IsFileExists(WEB_CONFIG_FILE_PATH) {
		content, err := os.ReadFile(WEB_CONFIG_FILE_PATH)
		if err != nil {
			tipLabel.SetText("文件错误：请放置到正确位置")
		} else {
			fmt.Printf("Read config: %s\n", string(content))
			json.Unmarshal(content, &webConfig)
			tipLabel.SetText(fmt.Sprintf("端口: %d", webConfig.HttpPort))
		}
	}

	openBrowser := widget.NewButton("访问面板", func() {
		if err := utils.Open(fmt.Sprintf("http://localhost:%d/", webConfig.HttpPort)); err != nil {
			fmt.Printf("Open Browser err %v\n", err)
		}
	})

	//守护进程管理
	pwd, _ := os.Getwd()
	fmt.Println("CWD:", pwd)
	// 程序所在目录
	daemon := cmd.NewProcessMgr(pwd+"/mcsmanager/daemon/", "./node_app.exe", "app.js")
	web := cmd.NewProcessMgr(pwd+"/mcsmanager/web/", "./node_app.exe", "app.js")

	//监听程序运行状态
	daemon.ListenStop(func(err error) {
		fmt.Println("EVENT: daemon exit!")
		if web.Started {
			web.End()
		}
		operationButton.Enable()
		operationButton.SetText("启动后台程序")
		statusLabel.SetText(STOPPED_TEXT)
		statusLabel.SetColor(color.Black)
	})
	web.ListenStop(func(err error) {
		fmt.Println("EVENT: web exit!")
		if daemon.Started {
			daemon.End()
		}
	})

	// 启动/关闭按钮点击事件
	operationButton.OnTapped = func() {
		var err error
		if !daemon.Started {
			if err = daemon.Start(); err != nil {
				utils.WriteErrLog(fmt.Sprintf("Start daemon error:%s", err.Error()))
				return
			}
			if err = web.Start(); err != nil {
				utils.WriteErrLog(fmt.Sprintf("Start web error:%s", err.Error()))
				daemon.End()
				return
			}
			statusLabel.SetText(STARTED_TEXT)
			statusLabel.SetColor(utils.GREEN)
			operationButton.SetText("停止后台程序")
		} else {
			operationButton.SetText("停止中...")
			operationButton.Disable()
			if err = web.End(); err != nil {
				utils.WriteErrLog(fmt.Sprintf("Stop daemon error:%s", err.Error()))
				return
			}
		}
	}

	window.SetCloseIntercept(func() {
		dialog.ShowConfirm("警告", "确定要退出程序吗？", func(b bool) {
			if b {
				if daemon.Started {
					dialog.ShowInformation("错误", "您必须关闭后台程序才能关闭本窗口", window)
					return
				}
				os.Exit(0)
			}
		}, window)
	})

	paddingContainer1 := container.New(layout.NewPaddedLayout(), infoLabel.Canvas)
	paddingContainer2 := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(), statusLabel.Canvas, tipLabelWrapper))
	paddingContainer3 := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(), btnWrapper, openBrowser))
	content := container.New(layout.NewVBoxLayout(), paddingContainer1, layout.NewSpacer(), paddingContainer2, paddingContainer3)
	window.SetContent(container.New(layout.NewPaddedLayout(), content))
	window.ShowAndRun()
}
