package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
)

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(1300, 800)
	window.SetWindowTitle("bend-it playground")

	NewPlayground(window)

	window.Show()

	app.Exec()
}
