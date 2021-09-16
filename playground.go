package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/walpod/bend-it"
	"github.com/walpod/bend-it/cubic"
)

// TODO how to keep in sync with spline
type GraphicsSceneItems struct {
	Scene        *widgets.QGraphicsScene
	segmentItems []*widgets.QGraphicsPathItem
	vertexItems  []*widgets.QGraphicsEllipseItem // per knot-no
}

func (si *GraphicsSceneItems) SetSegmentItem(segmentNo int, path gui.QPainterPath_ITF, pen gui.QPen_ITF, brush gui.QBrush_ITF) {
	// append to slice if necessary
	if segmentNo >= len(si.segmentItems) {
		newCnt := segmentNo - len(si.segmentItems) + 1
		si.segmentItems = append(si.segmentItems, make([]*widgets.QGraphicsPathItem, newCnt)...)
	}
	// remove old pathItem if exists
	if si.segmentItems[segmentNo] != nil {
		si.Scene.RemoveItem(si.segmentItems[segmentNo])
	}
	// set item
	pathItem := si.Scene.AddPath(path, pen, brush)
	si.segmentItems[segmentNo] = pathItem
}

func (si *GraphicsSceneItems) SegmentItem(segmentNo int) *widgets.QGraphicsPathItem {
	if segmentNo >= len(si.segmentItems) {
		return nil
	} else {
		return si.segmentItems[segmentNo]
	}
}

func (si *GraphicsSceneItems) SetVertexItem(knotNo int, vertexItem *widgets.QGraphicsEllipseItem) {
	// append to slice if necessary
	if knotNo >= len(si.vertexItems) {
		newCnt := knotNo - len(si.vertexItems) + 1
		si.vertexItems = append(si.vertexItems, make([]*widgets.QGraphicsEllipseItem, newCnt)...)
	}
	// remove old vertexItem if exists
	if si.vertexItems[knotNo] != nil {
		si.Scene.RemoveItem(si.vertexItems[knotNo])
	}
	// set item
	si.vertexItems[knotNo] = vertexItem
	si.Scene.AddItem(vertexItem)
}

func (si *GraphicsSceneItems) VertexItem(knotNo int) *widgets.QGraphicsEllipseItem {
	if knotNo >= len(si.vertexItems) {
		return nil
	} else {
		return si.vertexItems[knotNo]
	}
}

type Playground struct {
	spline     bendit.Spline2d
	sceneItems GraphicsSceneItems
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

	pg.sceneItems.Scene = widgets.NewQGraphicsScene(central)
	pg.sceneItems.Scene.SetSceneRect2(0, 0, float64(mainWindow.Width()), float64(mainWindow.Height()))

	view := widgets.NewQGraphicsView(central)
	view.SetScene(pg.sceneItems.Scene)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(view, 0, 0)

	central.SetLayout(layout)

	pg.buildSpline()
	pg.addSplineToScene()
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
	/*pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(200, 200, nil, cubic.NewControl(210, 200)),
		//cubic.NewBezierVx2(400, 400, cubic.NewControl(390, 400), nil)) TODO fix dependents
		cubic.NewBezierVx2(400, 400, cubic.NewControl(390, 400), cubic.NewControl(410, 400)))
	pg.spline = cubic.NewBezierSpline2d(nil,
		//cubic.NewBezierVx2(0, 0, nil, cubic.NewControl(100, 0)),
		cubic.NewBezierVx2(0, 0, cubic.NewControl(-100, 0), cubic.NewControl(100, 0)),
		//cubic.NewBezierVx2(100, 100, cubic.NewControl(0, 100), nil))
		cubic.NewBezierVx2(100, 100, cubic.NewControl(0, 100), cubic.NewControl(200, 100)))*/
	pg.spline = cubic.NewBezierSpline2d(nil,
		cubic.NewBezierVx2(100, 100, nil, cubic.NewControl(120, 150)),
		//cubic.NewBezierVx2(300, 300, cubic.NewControl(200, 300), nil), //cubic.NewControl(400, 300)),  TODO fix dependents : refresh dependent control ...
		cubic.NewBezierVx2(300, 300, cubic.NewControl(200, 300), cubic.NewControl(400, 300)),
		//cubic.NewBezierVx2(500, 100, cubic.NewControl(490, 150), nil))
		cubic.NewBezierVx2(500, 100, cubic.NewControl(490, 150), cubic.NewControl(510, 50)))
}

func (pg *Playground) addSplineToScene() {
	// colors
	black := gui.NewQColor2(core.Qt__black)
	gray := gui.NewQColor2(core.Qt__gray)

	// styles for spline and vertices
	//pen := gui.NewQPen3(gui.NewQColor3(0, 0, 0, 255))
	pen := gui.NewQPen3(black)
	brush := gui.NewQBrush2(core.Qt__SolidPattern)
	radius := 6.0

	// styles for controls
	//penCtrl := gui.NewQPen3(gray) //gui.NewQPen2(core.Qt__DotLine)
	brushCtrl := gui.NewQBrush3(gray, core.Qt__SolidPattern)

	// vertex as solid black circle
	addVertexToScene := func(knotNo int, x float64, y float64) {
		veh := VertexEventHandler{spline: pg.spline, knotNo: knotNo}
		ellipseVi := widgets.NewQGraphicsEllipseItem3(x-radius, y-radius, 2*radius, 2*radius, nil)
		ellipseVi.SetBrush(brush)
		ellipseVi.ConnectMousePressEvent(veh.HandleMousePressEvent)
		ellipseVi.ConnectMouseReleaseEvent(veh.HandleMouseReleaseEvent)
		//pg.sceneItems.Scene.AddItem(ellipseVi)
		pg.sceneItems.SetVertexItem(knotNo, ellipseVi)
	}

	// bezier-control as solid gray circle
	addBezierControlToScene := func(knotNo int, ctrl *cubic.Control, isEntry bool) {
		//scene.AddLine2(x, y, bvx.Entry().X(), bvx.Entry().Y(), penCt)
		ellCtrl := widgets.NewQGraphicsEllipseItem3(ctrl.X()-radius, ctrl.Y()-radius, 2*radius, 2*radius, nil)
		evh := BezierControlEventHandler{bezier: pg.spline.(*cubic.BezierSpline2d), knotNo: knotNo, isEntry: isEntry,
			radius: radius, playground: pg, ellCtrl: ellCtrl}
		ellCtrl.SetBrush(brushCtrl)
		ellCtrl.ConnectMousePressEvent(evh.HandleMousePressEvent)
		ellCtrl.ConnectMouseReleaseEvent(evh.HandleMouseReleaseEvent)
		pg.sceneItems.Scene.AddItem(ellCtrl)
	}

	// vertices
	knots := pg.spline.Knots()
	for i := 0; i < knots.Cnt(); i++ {
		t, _ := knots.Knot(i)

		x, y := pg.spline.At(t)
		addVertexToScene(i, x, y)

		// controls
		switch spl := pg.spline.(type) {
		case *cubic.BezierSpline2d:
			// bezier control points
			bvx, _ := spl.Vertex(i).(*cubic.BezierVx2)
			if i > 0 {
				addBezierControlToScene(i, bvx.Entry(), true)
			}
			if i < knots.Cnt()-1 {
				addBezierControlToScene(i, bvx.Exit(), false)
			}
		default:
			panic(fmt.Sprintf("type not yet supported: %T", spl))
		}
	}

	// line segments
	/*pg.spline.Approx(0.5, bendit.NewDirectCollector2d(func(tstart, tend, pstartx, pstarty, pendx, pendy float64) {
		scene.AddLine2(pstartx, pstarty, pendx, pendy, pen)
	}))*/
	paco := NewQPathCollector2d()
	pg.spline.Approx(0.5, paco)
	fmt.Printf("#line-segments: %v \n", paco.LineCnt())
	for segmNo, path := range paco.Paths {
		pg.sceneItems.SetSegmentItem(segmNo, path, pen, gui.NewQBrush())
	}
}

/*func (pg *Playground) paint(canvas *widgets.QWidget) {
	qp := gui.NewQPainter2(canvas)

	// draw spline
	//pg.drawSplineByIteration(qp)
	//pg.drawSplineBySubdivisionDirect(qp)
	//pg.drawSplineBySubdivisionPath(qp)

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
	collector := bendit.NewDirectCollector2d(func(segmentNo int, tstart, tend, pstartx, pstarty, pendx, pendy float64) {
		fmt.Printf("%v-th line(%v, %v, %v, %v)\n", lineSegNo, pstartx, pstarty, pendx, pendy)
		lineSegNo++
		qp.DrawLine3(int(math.Round(pstartx)), int(math.Round(pstarty)), int(math.Round(pendx)), int(math.Round(pendy)))
	})
	pg.spline.Approx(0.2, collector)
}*/

type VertexEventHandler struct {
	spline bendit.Spline2d
	knotNo int
	//ellVx  *widgets.QGraphicsEllipseItem
	//mousePressX, mousePressY float64
}

func (eh *VertexEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	//eh.mousePressX, eh.mousePressY = event.Pos().X(), event.Pos().Y()
	//fmt.Printf("mouse-press-event for vertex with knotNo = %v at %v/%v\n", eh.knotNo, eh.mousePressX, eh.mousePressY)
}

func (eh *VertexEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	vx := eh.spline.Vertex(eh.knotNo)
	knotX, knotY := vx.Coord()
	fmt.Printf("mouse-released-event for vertex with knotNo = %v at %v/%v, for knot previously at %v/%v\n",
		eh.knotNo, pos.X(), pos.Y(), knotX, knotY)

	// TODO move vertex
}

type BezierControlEventHandler struct {
	bezier     *cubic.BezierSpline2d
	knotNo     int
	isEntry    bool
	radius     float64
	ellCtrl    *widgets.QGraphicsEllipseItem
	playground *Playground
}

func (eh *BezierControlEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
}

func (eh *BezierControlEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	bvx := eh.bezier.BezierVertex(eh.knotNo)

	pos := event.Pos()
	ctrl := cubic.NewControl(pos.X(), pos.Y())

	// prepare new entry and exit controls
	var entry, exit *cubic.Control
	if eh.isEntry {
		entry = ctrl
		if !bvx.Dependent() {
			exit = bvx.Exit()
		}
	} else {
		exit = ctrl
		if !bvx.Dependent() {
			entry = bvx.Entry()
		}
	}

	// modify bezier
	x, y := bvx.Coord()
	eh.bezier.Update(eh.knotNo, x, y, entry, exit)

	// move control circles
	eh.ellCtrl.SetRect2(ctrl.X()-eh.radius, ctrl.Y()-eh.radius, 2*eh.radius, 2*eh.radius)

	// TODO redraw spline (at least one segment or two if dependent)
	if eh.isEntry || (eh.bezier.BezierVertex(eh.knotNo).Dependent() && eh.knotNo > 0) {
		segmNo := eh.knotNo - 1
		//qPathItem := eh.playground.sceneItems.SegmentItem(segmNo)
		//eh.playground.sceneItems.Scene.RemoveItem(qPathItem)
		paco := NewQPathCollector2d()
		eh.playground.spline.ApproxSegments(segmNo, segmNo, 0.5, paco)
		path := paco.Paths[segmNo]
		pen := gui.NewQPen3(gui.NewQColor2(core.Qt__black))
		//qPathItem = eh.playground.sceneItems.Scene.AddPath(path, pen, gui.NewQBrush())
		//eh.playground.sceneItems.SetSegmentItem(segmNo, qPathItem)
		eh.playground.sceneItems.SetSegmentItem(segmNo, path, pen, gui.NewQBrush())
	}

	if !eh.isEntry || (eh.bezier.BezierVertex(eh.knotNo).Dependent() && eh.knotNo < eh.bezier.Knots().Cnt()) {
		segmNo := eh.knotNo
		//qPathItem := eh.playground.sceneItems.SegmentItem(segmNo)
		//eh.playground.sceneItems.Scene.RemoveItem(qPathItem)
		paco := NewQPathCollector2d()
		eh.playground.spline.ApproxSegments(segmNo, segmNo, 0.5, paco)
		path := paco.Paths[segmNo]
		pen := gui.NewQPen3(gui.NewQColor2(core.Qt__black))
		//qPathItem = eh.playground.sceneItems.Scene.AddPath(path, pen, gui.NewQBrush())
		//eh.playground.sceneItems.SetSegmentItem(segmNo, qPathItem)
		eh.playground.sceneItems.SetSegmentItem(segmNo, path, pen, gui.NewQBrush())
	}

	/*lastx, lasty := bvx.Coord()
	fmt.Printf("mouse-released-event for vertex with knotNo = %vx at %vx/%vx, for knot previously at %vx/%vx\n",
		eh.knotNo, pos.X(), pos.Y(), lastx, lasty)*/
}

func (pg *Playground) drawSplineBySubdivisionPath(qp *gui.QPainter) {
	paco := NewQPathCollector2d()
	pg.spline.Approx(0.5, paco)
	fmt.Printf("#line-segments: %v \n", paco.LineCnt())
	pen := gui.NewQPen()
	for _, path := range paco.Paths {
		qp.StrokePath(path, pen)
	}
}

type QPathCollector2d struct {
	Paths map[int]*gui.QPainterPath
}

func NewQPathCollector2d() *QPathCollector2d {
	return &QPathCollector2d{Paths: map[int]*gui.QPainterPath{}}
}

func (lc *QPathCollector2d) CollectLine(segmentNo int, tstart, tend, pstartx, pstarty, pendx, pendy float64) {
	// get path for segment
	path, exists := lc.Paths[segmentNo]
	if !exists {
		path = gui.NewQPainterPath()
		lc.Paths[segmentNo] = path
	}

	// add line to path
	if path.ElementCount() == 0 {
		path.MoveTo(core.NewQPointF3(pstartx, pstarty))
	}
	path.LineTo(core.NewQPointF3(pendx, pendy))
}

func (lc *QPathCollector2d) LineCnt() int {
	lineCnt := 0
	for _, path := range lc.Paths {
		lineCnt += path.ElementCount()
	}
	return lineCnt
}

/*func (pg *Playground) drawTest(qp *gui.QPainter) {
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
}*/
