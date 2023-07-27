package main

import (
	"fmt"

	"github.com/MCSManager/Launcher/lang"
	"github.com/fatih/color"
)

func logErr(err string) {
	fmt.Println(color.HiRedString(lang.T("ERROR") + " " + err))
}

func logInfo(text string) {
	fmt.Println(color.HiGreenString(lang.T("INFO")) + " " + text)
}
