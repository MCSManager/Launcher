package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

//go build -ldflags -H=windowsgui main.go

// StartCmd 运行脚本并输出到 output 控件
func StartCmd(name string, outputW *widget.Label, args ...string) {
	args = append([]string{name}, args...)
	baseCmd := "sh"
	if runtime.GOOS == "windows" {
		//todo 待验证
		baseCmd = "cmd"
	}
	outputByt, err := exec.Command(baseCmd, args...).Output()
	if err != nil {
		outStr := fmt.Sprintf("%s 启动失败, error: %s", name, err.Error())
		//todo 考虑是否将 output 输出到控件上（可滚动控件）
		outputW.SetText(outStr)
		return
	}
	outputW.SetText(string(outputByt))
}

func main() {
	fontPath := "./config/msyh.ttc"
	os.Setenv("FYNE_FONT", fontPath)
	//fmt.Println("U %v", utils.IsFileExists(fontPath))
	a := app.New()
	w := a.NewWindow("MCSManager 面板管理小工具")

	w.Resize(fyne.Size{Width: 280, Height: 360})
	daemonLabel := widget.NewLabel("daemon output")
	daemonOutput := widget.NewLabel("")
	webLabel := widget.NewLabel("web output")
	webOutput := widget.NewLabel("")

	finishLabel := widget.NewLabel("")

	btn := widget.NewButton("启动", nil)
	btnToggle := false
	btn.OnTapped = func() {
		btnToggle = !btnToggle
		btn.Disabled()
		if btnToggle { //启动
			btn.SetText("启动中...")
			wg := sync.WaitGroup{}
			wg.Add(2)
			//启动程序1
			go func() {
				StartCmd("./out/test1", daemonOutput, "1")
				wg.Done()
			}()

			//启动程序2
			go func() {
				StartCmd("./out/test2", webOutput, "2")
				wg.Done()
			}()

			wg.Wait()
			btn.SetText("停止")
		} else { //停止
			btn.SetText("停止中...")
			time.Sleep(time.Second * 2)
			btn.SetText("启动")
		}
		finishLabel.SetText("完成")
		btn.Enable()
	}

	btn_color := canvas.NewRectangle(
		color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	container1 := container.New(
		// layout of container
		layout.NewMaxLayout(),
		// first use btn color
		btn_color,
		// 2nd btn widget
		btn,
	)

	content := container.New(layout.NewVBoxLayout(), daemonLabel, daemonOutput, webLabel, webOutput, layout.NewSpacer(), finishLabel, container1)

	w.SetContent(content)

	w.ShowAndRun()
}
