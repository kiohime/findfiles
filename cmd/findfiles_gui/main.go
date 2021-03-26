package main

import (
	"fmt"
	"os"

	"github.com/kiohime/findfiles/pkg/ffgui"
)

func main() {
	ffgui.AppSet = &ffgui.Settings{}
	ffgui.AppData = &ffgui.Data{}
	err := ffgui.Initialize(ffgui.AppSet)
	if err != nil {
		fmt.Printf("Error in initialisation : %q\n", err)
		ffgui.KeyWait()
		os.Exit(1)
	}
	ffgui.Gui(ffgui.AppSet, ffgui.AppData)
}
