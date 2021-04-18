package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/walpod/bend-it"
	"github.com/walpod/bend-it/cubic"
	"math"
)

type Playground struct {
	//window *widgets.QMainWindow
	//canvas *widgets.QWidget
	spline bendit.Spline2d
}

func (pg *Playground) build(window *widgets.QMainWindow) {
	//pg.window = window

	/*statusbar := widgets.NewQStatusBar(window)
	window.SetStatusBar(statusbar)
	statusbar.ShowMessage("the status bar ...", 0)*/

	central := widgets.NewQWidget(window, 0)
	central.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(central)

	canvas := widgets.NewQWidget(central, 0)
	//canvas.Resize2(800, 500)
	central.Layout().AddWidget(canvas)

	pg.buildSpline()

	canvas.ConnectPaintEvent(func(vqp *gui.QPaintEvent) {
		pg.paint(canvas)
	})
}

func (pg *Playground) buildSpline() {
	// hermite
	pg.spline = cubic.NewHermiteSplineTanFinder2d([]float64{10, 100, 150}, []float64{10, 100, 10}, cubic.NaturalTanf2d{}, nil)

	// canonical
	/*cubics := []cubic.Cubic2d{cubic.NewCubic2d(
		cubic.NewCubicPoly(100, 80, 40, 8),
		cubic.NewCubicPoly(210, 120, 0, 9),
	)}
	pg.spline = cubic.NewCanonicalSpline2d(cubics, nil)
	*/

	// bezier
	/*pg.spline = cubic.NewBezierSpline2d([]float64{200, 400}, []float64{200, 400}, []float64{210, 390}, []float64{200, 400}, nil)
	pg.spline = cubic.NewBezierSpline2d(
		[]float64{100, 300, 500}, []float64{100, 300, 100},
		[]float64{120, 230, 370, 490}, []float64{150, 300, 300, 150}, nil)
	*/
}

func (pg *Playground) paint(canvas *widgets.QWidget) {
	qp := gui.NewQPainter2(canvas)
	dom := pg.spline.Domain()
	stepSize := dom.To / 100
	for t := dom.From; t < dom.To; t += stepSize {
		x, y := pg.spline.At(t)
		qp.DrawPoint3(int(math.Round(x)), int(math.Round(y)))
	}
	qp.DestroyQPainter()
}
