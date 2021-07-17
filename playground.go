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
	spline bendit.Spline2d
}

/*func (pg *Playground) build(mainWindow *widgets.QMainWindow) {
	//pg.mainWindow = mainWindow

	//statusbar := widgets.NewQStatusBar(mainWindow)
	//mainWindow.SetStatusBar(statusbar)
	//statusbar.ShowMessage("the status bar ...", 0)

	central := widgets.NewQWidget(mainWindow, 0)
	central.SetLayout(widgets.NewQVBoxLayout())
	mainWindow.SetCentralWidget(central)

	canvas := widgets.NewQWidget(central, 0)
	canvas.Resize2(800, 500)
	central.Layout().AddWidget(canvas)

	pg.buildSpline()

	central.ConnectPaintEvent(func(vqp *gui.QPaintEvent) {
		pg.paint(central)
	})

	canvas.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		fmt.Printf("mouse event but no info about position (?): %v", *event)
	})
}*/

func (pg *Playground) buildWthScene(mainWindow *widgets.QMainWindow) {
	central := widgets.NewQWidget(mainWindow, 0)
	mainWindow.SetCentralWidget(central)

	scene := widgets.NewQGraphicsScene(central)
	scene.SetSceneRect2(0, 0, float64(mainWindow.Width()), float64(mainWindow.Height()))

	view := widgets.NewQGraphicsView(central)
	view.SetScene(scene)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(view, 0, 0)

	central.SetLayout(layout)

	pg.buildSpline()
	pg.addSplineToScene(scene)
}

func (pg *Playground) buildSpline() {
	// hermite
	/*pg.spline = cubic.NewNaturalHermiteSpline2d(nil,
		cubic.NewHermiteVx2Raw(10, 10),
		cubic.NewHermiteVx2Raw(100, 100),
		cubic.NewHermiteVx2Raw(150, 10),
	)
	pg.spline = cubic.NewNaturalHermiteSpline2d(nil,
		cubic.NewHermiteVx2Raw(100, 100),
		cubic.NewHermiteVx2Raw(400, 400),
		cubic.NewHermiteVx2Raw(700, 100),
	)*/
	herm := cubic.NewNaturalHermiteSpline2d(nil)
	herm.Add(cubic.NewHermiteVx2Raw(100, 100))
	herm.Add(cubic.NewHermiteVx2Raw(400, 400))
	herm.Add(cubic.NewHermiteVx2Raw(700, 100))
	//herm.InsCoord(1, 400, 400)
	pg.spline = herm

	// canonical
	/*pg.spline = cubic.NewCanonicalSpline2d(nil,
		cubic.NewCubic2d(cubic.NewCubicPoly(100, 80, 40, 8), cubic.NewCubicPoly(210, 120, 0, 9)),
	)*/

	// bezier
	/*pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(200, 200, 0, 0, 210, 200),
		cubic.NewBezierVx2(400, 400, 390, 400, 0, 0))
	pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(0, 0, 0, 0, 100, 0),
		cubic.NewBezierVx2(100, 100, 0, 100, 0, 0))*/
	/*pg.spline = cubic.NewBezierSpline2d(nil,
	cubic.NewBezierVx2(100, 100, 0, 0, 120, 150),
	cubic.NewBezierVx2(300, 300, 200, 300, 400, 300),
	cubic.NewBezierVx2(500, 100, 490, 150, 0, 0))*/
}

type VertexEventHandler struct {
	spline                   bendit.Spline2d
	knotNo                   int
	mousePressX, mousePressY float64
}

func (vc *VertexEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	//vc.mousePressX, vc.mousePressY = event.Pos().X(), event.Pos().Y()
	vc.mousePressX = event.Pos().X()
	fmt.Printf("mouse-press-event for vertex with knotNo = %v at %v/%v\n", vc.knotNo, vc.mousePressX, vc.mousePressY)
}

func (vc *VertexEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	if v, err := vc.spline.Vertex(vc.knotNo); err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		lastx, lasty := v.Coord()
		fmt.Printf("mouse-released-event for vertex with knotNo = %v at %v/%v, pressed at %v/%v, for knot previously at %v/%v\n",
			vc.knotNo, pos.X(), pos.Y(), vc.mousePressX, vc.mousePressY, lastx, lasty)

		// TODO move vertex
	}
}

func (pg Playground) addSplineToScene(scene *widgets.QGraphicsScene) {
	//red := gui.NewQColor3(255, 0, 0, 255)
	pen := gui.NewQPen3(gui.NewQColor3(0, 0, 0, 255))
	brush := gui.NewQBrush2(core.Qt__SolidPattern)
	radius := 6.0

	// vertices
	knots := pg.spline.Knots()
	for i := 0; i < knots.Count(); i++ {
		t, _ := knots.Knot(i)
		x, y := pg.spline.At(t)
		veh := VertexEventHandler{spline: pg.spline, knotNo: i}
		ell := widgets.NewQGraphicsEllipseItem3(x-radius, y-radius, 2*radius, 2*radius, nil)
		ell.SetBrush(brush)
		ell.ConnectMousePressEvent(veh.HandleMousePressEvent)
		/*ell.ConnectMouseMoveEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
			pos, lastPos := event.Pos(), event.LastPos()
			fmt.Printf("mouse moved to %v/%v, previously at %v/%v\n", pos.X(), pos.Y(), lastPos.X(), lastPos.Y())
		})*/
		ell.ConnectMouseReleaseEvent(veh.HandleMouseReleaseEvent)

		scene.AddItem(ell)
	}

	// spline
	/*pg.spline.Approx(0.5, bendit.NewDirectCollector2d(func(tstart, tend, pstartx, pstarty, pendx, pendy float64) {
		scene.AddLine2(pstartx, pstarty, pendx, pendy, pen)
	}))*/

	paco := NewQPathCollector2d()
	pg.spline.Approx(0.5, paco)
	fmt.Printf("#line-segments: %v \n", paco.Path.ElementCount())
	scene.AddPath(paco.Path, pen, gui.NewQBrush())
	path := widgets.NewQGraphicsPathItem2(paco.Path, nil)
	path.SetPen(pen)
	scene.AddItem(path)
	//fmt.Printf("scene items# = %v", len(scene.Items2(core.NewQPointF3(100,100), core.Qt__IntersectsItemShape, core.Qt__AscendingOrder, gui.NewQTransform2())))

	scene.RemoveItem(path)
}

func (pg *Playground) paint(canvas *widgets.QWidget) {
	qp := gui.NewQPainter2(canvas)

	// draw spline
	//pg.drawSplineByIteration(qp)
	//pg.drawSplineBySubdivisionDirect(qp)
	pg.drawSplineBySubdivisionPath(qp)

	// draw spline vertices
	qp.Brush().SetStyle(core.Qt__SolidPattern)
	knots := pg.spline.Knots()
	for i := 0; i < knots.Count(); i++ {
		t, _ := knots.Knot(i)
		x, y := pg.spline.At(t)
		qp.DrawEllipse4(core.NewQPointF3(x, y), 5, 5)
	}

	//pg.drawTest(qp)
	qp.DestroyQPainter()
}

func (pg *Playground) drawSplineByIteration(qp *gui.QPainter) {
	knots := pg.spline.Knots()
	tend := knots.Tend()
	stepSize := tend / 100
	for t := knots.Tstart(); t < tend; t += stepSize {
		x, y := pg.spline.At(t)
		qp.DrawPoint3(int(math.Round(x)), int(math.Round(y)))
	}
}

func (pg *Playground) drawSplineBySubdivisionDirect(qp *gui.QPainter) {
	// draw spline
	lineSegNo := 0
	collector := bendit.NewDirectCollector2d(func(tstart, tend, pstartx, pstarty, pendx, pendy float64) {
		fmt.Printf("%v-th line(%v, %v, %v, %v)\n", lineSegNo, pstartx, pstarty, pendx, pendy)
		lineSegNo++
		qp.DrawLine3(int(math.Round(pstartx)), int(math.Round(pstarty)), int(math.Round(pendx)), int(math.Round(pendy)))
	})
	pg.spline.Approx(0.2, collector)
}

type QPathCollector2d struct {
	Path *gui.QPainterPath
}

func NewQPathCollector2d() *QPathCollector2d {
	return &QPathCollector2d{Path: gui.NewQPainterPath()}
}

func (lc QPathCollector2d) CollectLine(tstart, tend, pstartx, pstarty, pendx, pendy float64) {
	if lc.Path.ElementCount() == 0 {
		lc.Path.MoveTo(core.NewQPointF3(pstartx, pstarty))
	}
	lc.Path.LineTo(core.NewQPointF3(pendx, pendy))
}

func (pg *Playground) drawSplineBySubdivisionPath(qp *gui.QPainter) {
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
