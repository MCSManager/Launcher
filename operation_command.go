package main

import (
	"github.com/MCSManager/Launcher/lang"
	"github.com/rivo/tview"
)

type OperationCommand struct {
	MainText      string
	SecondaryText string
	Shortcut      rune
	Exec          func()
}

var globalOperationCommand = []OperationCommand{}

func initOperationCommand() {
	globalOperationCommand = []OperationCommand{
		{MainText: lang.T("TopCommand1"), SecondaryText: lang.T("TopCommandSubTitle1"), Shortcut: '1', Exec: func() {

		}},
		{MainText: lang.T("TopCommand2"), SecondaryText: lang.T("TopCommandSubTitle2"), Shortcut: '2', Exec: func() {

		}},
		{MainText: lang.T("TopCommand3"), SecondaryText: lang.T("TopCommandSubTitle3"), Shortcut: '3', Exec: func() {

		}},
		{MainText: lang.T("TopCommand4"), SecondaryText: lang.T("TopCommandSubTitle4"), Shortcut: '4', Exec: func() {

		}},
		{MainText: lang.T("TopCommand5"), SecondaryText: lang.T("TopCommandSubTitle5"), Shortcut: '5', Exec: func() {

		}},
		{MainText: lang.T("TopCommand6"), SecondaryText: lang.T("TopCommandSubTitle6"), Shortcut: '6', Exec: func() {

		}},
	}

}

type UIBox interface {
	SetBorder(show bool) *tview.Box
	SetTitle(title string) *tview.Box
	SetTitleAlign(align int) *tview.Box
	SetBorderPadding(left, top, right, bottom int) *tview.Box
}

func initUIBox(box UIBox, title string) {
	box.SetBorder(true)
	box.SetBorderPadding(0, 0, 1, 1)
	box.SetTitleAlign(tview.AlignLeft)
	box.SetTitle(title)
}
