package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
)

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(800, 500)
	window.SetWindowTitle("bend-it playground")

	var playground Playground
	playground.build(window)

	window.Show()

	app.Exec()
}
