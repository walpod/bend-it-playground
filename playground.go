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
	pg.spline = cubic.NewCardinalHermiteSpline2d(
		bendit.NewUniformKnots(), 0,
		cubic.NewRawHermiteVertex2d(100, 100),
		cubic.NewRawHermiteVertex2d(400, 400),
		cubic.NewRawHermiteVertex2d(700, 100),
	)

	// canonical
	/*pg.spline = cubic.NewCanonicalSpline2d(bendit.NewUniformKnots(),
		cubic.NewCubic2d(cubic.NewCubicPoly(100, 80, 40, 8), cubic.NewCubicPoly(210, 120, 0, 9)),
	)*/

	// bezier
	/*pg.spline = cubic.NewBezierSpline2d(bendit.NewUniformKnots(),
		cubic.NewBezierVertex2d(200, 200, 0, 0, 210, 200),
		cubic.NewBezierVertex2d(400, 400, 390, 400, 0, 0))
	pg.spline = cubic.NewBezierSpline2d(
		bendit.NewUniformKnots(),
		cubic.NewBezierVertex2d(0, 0, 0, 0, 100, 0),
		cubic.NewBezierVertex2d(100, 100, 0, 100, 0, 0))
	*/
	/*pg.spline = cubic.NewBezierSpline2d(bendit.NewUniformKnots(),
	cubic.NewBezierVertex2d(100, 100, 0, 0, 120, 150),
	cubic.NewBezierVertex2d(300, 300, 200, 300, 400, 300),
	cubic.NewBezierVertex2d(500, 100, 490, 150, 0, 0))
	*/
}

func (pg *Playground) paint(canvas *widgets.QWidget) {
	qp := gui.NewQPainter2(canvas)
	pg.drawByIteration(qp)
	//pg.drawBySubdivisionDirect(qp)
	//pg.drawBySubdivisionPath(qp)
	//pg.drawTest(qp)
	qp.DestroyQPainter()
}

func (pg *Playground) drawByIteration(qp *gui.QPainter) {
	dom := pg.spline.Knots().Domain(pg.spline.SegmentCnt())
	stepSize := dom.End / 100
	for t := dom.Start; t < dom.End; t += stepSize {
		x, y := pg.spline.At(t)
		qp.DrawPoint3(int(math.Round(x)), int(math.Round(y)))
	}
}

func (pg *Playground) drawBySubdivisionDirect(qp *gui.QPainter) {
	lineSegNo := 0
	collector := bendit.NewDirectCollector2d(func(ts, te, sx, sy, ex, ey float64) {
		fmt.Printf("%v-th line(%v, %v, %v, %v)\n", lineSegNo, sx, sy, ex, ey)
		lineSegNo++
		qp.DrawLine3(int(math.Round(sx)), int(math.Round(sy)), int(math.Round(ex)), int(math.Round(ey)))
	})
	pg.spline.Approx(0.2, collector)
}

type QPathCollector2d struct {
	Path *gui.QPainterPath
}

func NewQPathCollector2d() *QPathCollector2d {
	return &QPathCollector2d{Path: gui.NewQPainterPath()}
}

func (lc QPathCollector2d) CollectLine(ts, te, sx, sy, ex, ey float64) {
	if lc.Path.ElementCount() == 0 {
		lc.Path.MoveTo(core.NewQPointF3(sx, sy))
	}
	lc.Path.LineTo(core.NewQPointF3(ex, ey))
}

func (pg *Playground) drawBySubdivisionPath(qp *gui.QPainter) {
	paco := NewQPathCollector2d()
	pg.spline.Approx(0.5, paco)
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
