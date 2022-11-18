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
	"mcsmanager.com/desktop-app/utils"
	"os"
)

//go build -ldflags -H=windowsgui main.go

func main() {

	os.Setenv("FYNE_FONT", "C:/Windows/Fonts/msyhl.ttc")
	fmt.Println("U %v", utils.IsFileExists("C:/Windows/Fonts/msyhl.ttc"))
	a := app.New()
	w := a.NewWindow("MCSManager 面板管理小工具")

	w.Resize(fyne.Size{Width: 280, Height: 360})
	hello := widget.NewLabel("Hello 这是一个测试程序，请点按钮!")
	hello3 := widget.NewLabel("Hello 66666666666")
	hello2 := widget.NewLabel("曹操!")

	btn := widget.NewButton("点一下哦~", func() {
		hello.SetText("Welcome :)")
	})
	btn_color := canvas.NewRectangle(
		color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	container1 := container.New(
		// layout of container
		layout.NewMaxLayout(),
		// first use btn color
		btn_color,
		// 2nd btn widget
		btn,
	)






	content := container.New(layout.NewVBoxLayout(), hello, hello2, layout.NewSpacer(),hello3, container1)

	w.SetContent(content)

	w.ShowAndRun()
}
