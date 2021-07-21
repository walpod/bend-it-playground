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
	)
	herm := cubic.NewNaturalHermiteSpline2d(nil)
	herm.Add(cubic.NewHermiteVx2Raw(100, 100))
	herm.Add(cubic.NewHermiteVx2Raw(400, 400))
	herm.Add(cubic.NewHermiteVx2Raw(700, 100))
	//herm.InsCoord(1, 400, 400)
	pg.spline = herm
	*/

	// canonical
	/*pg.spline = cubic.NewCanonicalSpline2d(nil,
		cubic.NewCubic2d(cubic.NewCubicPoly(100, 80, 40, 8), cubic.NewCubicPoly(210, 120, 0, 9)),
	)*/

	// bezier
	pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(200, 200, nil, cubic.NewControl(210, 200)),
		cubic.NewBezierVx2(400, 400, cubic.NewControl(390, 400), nil))
	pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(0, 0, nil, cubic.NewControl(100, 0)),
		cubic.NewBezierVx2(100, 100, cubic.NewControl(0, 100), nil))
	pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(100, 100, nil, cubic.NewControl(120, 150)),
		cubic.NewBezierVx2(300, 300, cubic.NewControl(200, 300), cubic.NewReflective()), //cubic.NewControl(400, 300)),
		cubic.NewBezierVx2(500, 100, cubic.NewControl(490, 150), nil))
}

type AnnexToVertex struct {
	ellVt    *widgets.QGraphicsEllipseItem
	ellEntry *widgets.QGraphicsEllipseItem
	ellExit  *widgets.QGraphicsEllipseItem
}

type VertexEventHandler struct {
	spline bendit.Spline2d
	knotNo int
	//mousePressX, mousePressY float64
}

func (eh *VertexEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	//eh.mousePressX, eh.mousePressY = event.Pos().X(), event.Pos().Y()
	//fmt.Printf("mouse-press-event for vertex with knotNo = %v at %v/%v\n", eh.knotNo, eh.mousePressX, eh.mousePressY)
}

func (eh *VertexEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	vt := eh.spline.Vertex(eh.knotNo)
	knotX, knotY := vt.Coord()
	fmt.Printf("mouse-released-event for vertex with knotNo = %vt at %vt/%vt, for knot previously at %vt/%vt\n",
		eh.knotNo, pos.X(), pos.Y(), knotX, knotY)

	// TODO move vertex
}

type BezierControlEventHandler struct {
	bezier  *cubic.BezierSpline2d
	knotNo  int
	isEntry bool
	radius  float64
}

func (eh *BezierControlEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
}

func (eh *BezierControlEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	bvt := eh.bezier.BezierVertex(eh.knotNo)
	entry := bvt.Entry()
	exit := bvt.Exit()
	annexVt := eh.bezier.Annex().GetFromVertex(eh.knotNo).(*AnnexToVertex)

	pos := event.Pos()
	ct := cubic.NewControl(pos.X(), pos.Y())

	// entry or exit
	if eh.isEntry {
		if bvt.Entry().IsCalculated() {
			bvt.FlipCalculatedController()
		}
		entry, exit = ct, bvt.Exit()
	} else {
		if bvt.Exit().IsCalculated() {
			bvt.FlipCalculatedController()
		}
		entry, exit = bvt.Entry(), ct
	}

	// modify bezier
	x, y := bvt.Coord()
	eh.bezier.Update(eh.knotNo, x, y, entry, exit)

	// move control circles
	if entry != nil {
		annexVt.ellEntry.SetRect2(entry.ControlX()-eh.radius, entry.ControlY()-eh.radius, 2*eh.radius, 2*eh.radius)
	}
	if exit != nil {
		annexVt.ellExit.SetRect2(exit.ControlX()-eh.radius, exit.ControlY()-eh.radius, 2*eh.radius, 2*eh.radius)
	}

	// TODO redraw spline (at least part of it)

	/*lastx, lasty := bvt.Coord()
	fmt.Printf("mouse-released-event for vertex with knotNo = %vx at %vx/%vx, for knot previously at %vx/%vx\n",
		eh.knotNo, pos.X(), pos.Y(), lastx, lasty)*/
}

func (pg Playground) addSplineToScene(scene *widgets.QGraphicsScene) {
	// colors
	black := gui.NewQColor2(core.Qt__black)
	gray := gui.NewQColor2(core.Qt__gray)

	// styles for spline and vertices
	//pen := gui.NewQPen3(gui.NewQColor3(0, 0, 0, 255))
	pen := gui.NewQPen3(black)
	brush := gui.NewQBrush2(core.Qt__SolidPattern)
	radius := 6.0

	// styles for controls
	//penCt := gui.NewQPen3(gray) //gui.NewQPen2(core.Qt__DotLine)
	brushCt := gui.NewQBrush3(gray, core.Qt__SolidPattern)

	// vertex as solid black circle
	addVertexToScene := func(knotNo int, x float64, y float64, annexVt *AnnexToVertex) {
		ellVt := widgets.NewQGraphicsEllipseItem3(x-radius, y-radius, 2*radius, 2*radius, nil)
		veh := VertexEventHandler{spline: pg.spline, knotNo: knotNo}
		ellVt.SetBrush(brush)
		ellVt.ConnectMousePressEvent(veh.HandleMousePressEvent)
		ellVt.ConnectMouseReleaseEvent(veh.HandleMouseReleaseEvent)
		scene.AddItem(ellVt)
		annexVt.ellVt = ellVt
	}

	// bezier-control as solid gray circle
	addBezierControlToScene := func(knotNo int, ct cubic.Controller, isEntry bool, annexVt *AnnexToVertex) {
		if ct != nil {
			//scene.AddLine2(x, y, bvx.Entry().ControlX(), bvx.Entry().ControlY(), penCt)
			ellCt := widgets.NewQGraphicsEllipseItem3(ct.ControlX()-radius, ct.ControlY()-radius, 2*radius, 2*radius, nil)
			cteh := BezierControlEventHandler{bezier: pg.spline.(*cubic.BezierSpline2d), knotNo: knotNo, isEntry: isEntry, radius: radius}
			ellCt.SetBrush(brushCt)
			ellCt.ConnectMousePressEvent(cteh.HandleMousePressEvent)
			ellCt.ConnectMouseReleaseEvent(cteh.HandleMouseReleaseEvent)
			scene.AddItem(ellCt)
			if isEntry {
				annexVt.ellEntry = ellCt
			} else {
				annexVt.ellExit = ellCt
			}
		}
	}

	// vertices
	knots := pg.spline.Knots()
	for i := 0; i < knots.Count(); i++ {
		t, _ := knots.Knot(i)

		annexVt := &AnnexToVertex{}
		pg.spline.Annex().AttachToVertex(i, annexVt)

		x, y := pg.spline.At(t)
		addVertexToScene(i, x, y, annexVt)

		// controls
		switch spl := pg.spline.(type) {
		case *cubic.BezierSpline2d:
			// bezier control points
			bvx, _ := spl.Vertex(i).(*cubic.BezierVx2)
			addBezierControlToScene(i, bvx.Entry(), true, annexVt)
			addBezierControlToScene(i, bvx.Exit(), false, annexVt)
		}
	}

	// line segments
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
