package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/walpod/bend-it"
	"github.com/walpod/bend-it/cubic"
	"math"
)

type Playground struct {
	//window *widgets.QMainWindow
	//canvas *widgets.QWidget
	//spline *cubic.BezierSpline2d // TODO for a short time, then change back to bendit.Spline2d
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
	pg.spline = cubic.NewHermiteSplineTanFinder2d([]float64{10, 100, 150}, []float64{10, 100, 10}, cubic.NaturalTanf2d{}, bendit.NewUniformKnots())

	// canonical
	cubics := []cubic.Cubic2d{cubic.NewCubic2d(
		cubic.NewCubicPoly(100, 80, 40, 8),
		cubic.NewCubicPoly(210, 120, 0, 9),
	)}
	pg.spline = cubic.NewCanonicalSpline2d(cubics, bendit.NewUniformKnots())

	// bezier
	pg.spline = cubic.NewBezierSpline2d([]float64{200, 400}, []float64{200, 400},
		//[]float64{210, 390}, []float64{200, 400},
		[]float64{0, 390}, []float64{0, 400}, []float64{210, 0}, []float64{200, 0},
		bendit.NewUniformKnots())
	pg.spline = cubic.NewBezierSpline2d(
		[]float64{100, 300, 500}, []float64{100, 300, 100},
		//[]float64{120, 200, 400, 490}, []float64{150, 300, 300, 150},
		[]float64{0, 200, 490}, []float64{0, 300, 150}, []float64{120, 400, 0}, []float64{150, 300, 0},
		bendit.NewUniformKnots())
}

func (pg *Playground) paint(canvas *widgets.QWidget) {
	qp := gui.NewQPainter2(canvas)
	//pg.drawByIteration(qp)
	//pg.drawBySubdivisionDirect(qp)
	pg.drawBySubdivisionPath(qp)
	//pg.drawTest(qp)
	qp.DestroyQPainter()
}

func (pg *Playground) drawByIteration(qp *gui.QPainter) {
	dom := pg.spline.Knots().Domain(pg.spline.SegmentCnt())
	stepSize := dom.To / 100
	for t := dom.From; t < dom.To; t += stepSize {
		x, y := pg.spline.At(t)
		qp.DrawPoint3(int(math.Round(x)), int(math.Round(y)))
	}
}

func (pg *Playground) drawBySubdivisionDirect(qp *gui.QPainter) {
	lineSegNo := 0
	collector := cubic.NewDirectCollector2d(func(x0, y0, x1, y1 float64) {
		fmt.Printf("%v-th line(%v, %v, %v, %v)\n", lineSegNo, x0, y0, x1, y1)
		lineSegNo++
		qp.DrawLine3(int(math.Round(x0)), int(math.Round(y0)), int(math.Round(x1)), int(math.Round(y1)))
	})
	pg.spline.Approximate(0.2, collector)
}

type QPathCollector2d struct {
	Path *gui.QPainterPath
}

func NewQPathCollector2d() *QPathCollector2d {
	return &QPathCollector2d{Path: gui.NewQPainterPath()}
}

func (lc QPathCollector2d) CollectLine(x0, y0, x3, y3 float64) {
	if lc.Path.ElementCount() == 0 {
		lc.Path.MoveTo(core.NewQPointF3(x0, y0))
	}
	lc.Path.LineTo(core.NewQPointF3(x3, y3))
}

func (pg *Playground) drawBySubdivisionPath(qp *gui.QPainter) {
	paco := NewQPathCollector2d()
	pg.spline.Approximate(0.7, paco)
	fmt.Printf("#line-segments: %v \n", paco.Path.ElementCount())
	qp.StrokePath(paco.Path, gui.NewQPen())
}

func (pg *Playground) drawTest(qp *gui.QPainter) {
	pointsf := []*core.QPointF{
		core.NewQPointF3(250, 287),
		core.NewQPointF3(254, 289),
		core.NewQPointF3(259, 291),
		core.NewQPointF3(263, 293),
		core.NewQPointF3(500, 100),
	}
	points := [2]*core.QPoint{
		core.NewQPoint2(250, 287),
		core.NewQPoint2(500, 100),
	}
	_, _ = pointsf, points

	//qp.DrawPolyline(points[0], len(points))
	//qp.DrawPolyline(core.NewQPointF3(500, 200), 3)
	//qp.DrawPolygon3(points[0], len(points), core.Qt__OddEvenFill)
	//qp.DrawPoints(pointsf[0], len(points))

	path := gui.NewQPainterPath()
	path.MoveTo(core.NewQPointF3(100, 100))
	path.LineTo(core.NewQPointF3(250, 287))
	path.LineTo(core.NewQPointF3(259, 291))
	path.LineTo(core.NewQPointF3(263, 293))
	path.LineTo(core.NewQPointF3(500, 100))
	qp.StrokePath(path, gui.NewQPen())
}
