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
	scene          *widgets.QGraphicsScene
	segmentPaths   []*widgets.QGraphicsPathItem
	vertexCircles  []*widgets.QGraphicsEllipseItem    // per knot-no
	controlCircles [][2]*widgets.QGraphicsEllipseItem // entry and exit per knot-no
	controlLines   [][2]*widgets.QGraphicsLineItem    // entry and exit per knot-no
}

func NewGraphicsSceneItems(scene *widgets.QGraphicsScene) *GraphicsSceneItems {
	return &GraphicsSceneItems{scene: scene}
}

func (si *GraphicsSceneItems) SetSegmentPath(segmentNo int, path gui.QPainterPath_ITF,
	pen gui.QPen_ITF, brush gui.QBrush_ITF) *widgets.QGraphicsPathItem {

	// append to slice if necessary
	if segmentNo >= len(si.segmentPaths) {
		newCnt := segmentNo - len(si.segmentPaths) + 1
		si.segmentPaths = append(si.segmentPaths, make([]*widgets.QGraphicsPathItem, newCnt)...)
	}
	// remove old pathItem if exists
	if si.segmentPaths[segmentNo] != nil {
		si.scene.RemoveItem(si.segmentPaths[segmentNo])
	}
	// set item
	pathItem := si.scene.AddPath(path, pen, brush)
	si.segmentPaths[segmentNo] = pathItem
	return pathItem
}

func (si *GraphicsSceneItems) SetVertexCircle(knotNo int, circle *widgets.QGraphicsEllipseItem) {
	// append to slice if necessary
	if knotNo >= len(si.vertexCircles) {
		newCnt := knotNo - len(si.vertexCircles) + 1
		si.vertexCircles = append(si.vertexCircles, make([]*widgets.QGraphicsEllipseItem, newCnt)...)
	}
	// remove old circle if exists
	if si.vertexCircles[knotNo] != nil {
		si.scene.RemoveItem(si.vertexCircles[knotNo])
	}
	// set item
	si.vertexCircles[knotNo] = circle
	si.scene.AddItem(circle)
}

func (si *GraphicsSceneItems) VertexCircle(knotNo int) *widgets.QGraphicsEllipseItem {
	if knotNo >= len(si.vertexCircles) {
		return nil
	} else {
		return si.vertexCircles[knotNo]
	}
}

func (si *GraphicsSceneItems) mapToSideNo(isEntry bool) int {
	// map 'entry' to 0, 'exit' to 1
	if isEntry {
		return 0
	} else {
		return 1
	}
}

func (si *GraphicsSceneItems) SetControlCircle(knotNo int, isEntry bool, circle *widgets.QGraphicsEllipseItem) {
	// append to slice if necessary
	if knotNo >= len(si.controlCircles) {
		newCnt := knotNo - len(si.controlCircles) + 1
		si.controlCircles = append(si.controlCircles, make([][2]*widgets.QGraphicsEllipseItem, newCnt)...)
	}
	sideNo := si.mapToSideNo(isEntry)
	// remove old circle if exists
	if si.controlCircles[knotNo][sideNo] != nil {
		si.scene.RemoveItem(si.controlCircles[knotNo][sideNo])
	}
	// set item
	si.controlCircles[knotNo][sideNo] = circle
	si.scene.AddItem(circle)
}

func (si *GraphicsSceneItems) ControlCircle(knotNo int, isEntry bool) *widgets.QGraphicsEllipseItem {
	if knotNo >= len(si.controlCircles) {
		return nil
	} else {
		return si.controlCircles[knotNo][si.mapToSideNo(isEntry)]
	}
}

func (si *GraphicsSceneItems) SetControlLine(knotNo int, isEntry bool, fromx, fromy, tox, toy float64, pen gui.QPen_ITF) *widgets.QGraphicsLineItem {
	// append to slice if necessary
	if knotNo >= len(si.controlLines) {
		newCnt := knotNo - len(si.controlLines) + 1
		si.controlLines = append(si.controlLines, make([][2]*widgets.QGraphicsLineItem, newCnt)...)
	}
	sideNo := si.mapToSideNo(isEntry)
	// remove old line if exists
	if si.controlLines[knotNo][sideNo] != nil {
		si.scene.RemoveItem(si.controlLines[knotNo][sideNo])
	}
	// set line
	lineItem := si.scene.AddLine2(fromx, fromy, tox, toy, pen)
	si.controlLines[knotNo][sideNo] = lineItem
	return lineItem
}

/*func (si *GraphicsSceneItems) ResetVertex(knotNo int) {
	si.scene.RemoveItem(si.vertexCircles[knotNo])
	si.vertexCircles[knotNo] = nil
	si.scene.RemoveItem(si.controlCircles[knotNo][0])
	si.controlCircles[knotNo][0] = nil
	si.scene.RemoveItem(si.controlLines[knotNo][0])
	si.controlLines[knotNo][0] = nil
	si.scene.RemoveItem(si.controlCircles[knotNo][1])
	si.controlCircles[knotNo][1] = nil
	si.scene.RemoveItem(si.controlLines[knotNo][1])
	si.controlLines[knotNo][1] = nil
	prevSegmNo := knotNo - 1
	if prevSegmNo >= 0 && prevSegmNo < len(si.segmentPaths) {
		si.scene.RemoveItem(si.segmentPaths[prevSegmNo])
		si.segmentPaths[prevSegmNo] = nil
	}
	nextSegmNo := knotNo
	if nextSegmNo >= 0 && nextSegmNo < len(si.segmentPaths) {
		si.scene.RemoveItem(si.segmentPaths[nextSegmNo])
		si.segmentPaths[nextSegmNo] = nil
	}
}*/

/*func (si *GraphicsSceneItems) ControlLine(knotNo int, isEntry bool) *widgets.QGraphicsLineItem {
	if knotNo >= len(si.controlLines) {
		return nil
	} else {
		return si.controlLines[knotNo][si.mapToSideNo(isEntry)]
	}
}*/

type Playground struct {
	spline     bendit.VertSpline2d
	sceneItems GraphicsSceneItems
	// styles for spline and vertices
	pen   gui.QPen_ITF
	brush gui.QBrush_ITF
	// styles for controls
	penCtrl   gui.QPen_ITF
	brushCtrl gui.QBrush_ITF
}

func NewPlayground(mainWindow *widgets.QMainWindow) *Playground {
	central := widgets.NewQWidget(mainWindow, 0)
	mainWindow.SetCentralWidget(central)

	scene := widgets.NewQGraphicsScene(central)
	scene.SetSceneRect2(0, 0, float64(mainWindow.Width()), float64(mainWindow.Height()))

	view := widgets.NewQGraphicsView(central)
	view.SetScene(scene)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(view, 0, 0)

	central.SetLayout(layout)

	pg := &Playground{}
	pg.sceneItems = *NewGraphicsSceneItems(scene)

	// colors and styles
	black := gui.NewQColor2(core.Qt__black)
	gray := gui.NewQColor2(core.Qt__gray)
	pg.pen = gui.NewQPen3(black)
	pg.brush = gui.NewQBrush2(core.Qt__SolidPattern)
	pg.penCtrl = gui.NewQPen2(core.Qt__DotLine) // core.Qt__DashLine
	pg.brushCtrl = gui.NewQBrush3(gray, core.Qt__SolidPattern)

	pg.buildSpline()
	pg.addSplineToScene()
	return pg
}

/*func (pg *Playground) buildOld(mainWindow *widgets.QMainWindow) {
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

func (pg *Playground) buildSpline() {
	// hermite
	/*pg.spline = cubic.NewHermiteSpline2d(nil,
		cubic.NewHermiteVx2(200, 200, nil, cubic.NewControl(90, 90)),
		cubic.NewHermiteVx2(350, 350, cubic.NewControl(200, 0), nil),
		cubic.NewHermiteVx2(500, 200, cubic.NewControl(100, -100), nil),
	)*/

	/*pg.spline = cubic.NewNaturalHermiteSpline2d(nil,
		cubic.NewHermiteVx2Raw(10, 10),
		cubic.NewHermiteVx2Raw(100, 100),
		cubic.NewHermiteVx2Raw(150, 10),
	)*/

	/*pg.spline = cubic.NewNaturalHermiteSpline2d(nil,
		cubic.NewHermiteVx2Raw(100, 100),
		cubic.NewHermiteVx2Raw(400, 400),
		cubic.NewHermiteVx2Raw(700, 100),
	)*/

	/*herm := cubic.NewNaturalHermiteSpline2d(nil)
	herm.AddVertex(0, cubic.NewHermiteVx2Raw(100, 100))
	herm.AddVertex(1, cubic.NewHermiteVx2Raw(400, 400))
	herm.AddVertex(2, cubic.NewHermiteVx2Raw(700, 100))
	pg.spline = herm*/

	// bezier
	/*pg.spline = cubic.NewBezierSpline2d(nil,
	cubic.NewBezierVx2(200, 200, nil, cubic.NewControl(250, 200)),
	cubic.NewBezierVx2(400, 400, cubic.NewControl(350, 400), nil))*/

	/*pg.spline = cubic.NewBezierSpline2d(nil,
	cubic.NewBezierVx2(200, 200, cubic.NewControl(100, 200), cubic.NewControl(300, 200)),
	cubic.NewBezierVx2(300, 300, cubic.NewControl(200, 300), cubic.NewControl(400, 300)))*/

	pg.spline = cubic.NewBezierSpline2d(nil)
	pg.spline.AddVertex(0, cubic.NewBezierVx2(100, 100, nil, cubic.NewControl(120, 150)))
	pg.spline.AddVertex(1, cubic.NewBezierVx2(300, 300, cubic.NewControl(200, 300), nil))
	pg.spline.AddVertex(2, cubic.NewBezierVx2(500, 100, cubic.NewControl(490, 150), nil))
}

func (pg *Playground) vertexRectForCircle(x float64, y float64) *core.QRectF {
	radius := 6.0
	return core.NewQRectF4(x-radius, y-radius, 2*radius, 2*radius)
}

func (pg *Playground) addVertexToScene(knotNo int, x float64, y float64) {
	veh := VertexEventHandler{playground: pg, knotNo: knotNo}
	// vertex as solid black circle
	circleVt := widgets.NewQGraphicsEllipseItem2(pg.vertexRectForCircle(x, y), nil)
	circleVt.SetBrush(pg.brush)
	circleVt.ConnectMousePressEvent(veh.HandleMousePressEvent)
	circleVt.ConnectMouseReleaseEvent(veh.HandleMouseReleaseEvent)
	circleVt.ConnectMouseDoubleClickEvent(veh.HandleMouseDoubleClickEvent)
	pg.sceneItems.SetVertexCircle(knotNo, circleVt)
}

func (pg *Playground) controlRectForCircle(x float64, y float64) *core.QRectF {
	radius := 5.0
	return core.NewQRectF4(x-radius, y-radius, 2*radius, 2*radius)
}

func (pg *Playground) addControlPointToScene(knotNo int, vertex bendit.Vertex2d, ctrlx, ctrly float64, isEntry bool) {
	evh := ControlPointEventHandler{playground: pg, knotNo: knotNo, isEntry: isEntry}
	// bezier-control as solid gray circle
	circleCtrl := widgets.NewQGraphicsEllipseItem2(pg.controlRectForCircle(ctrlx, ctrly), nil)
	circleCtrl.SetBrush(pg.brushCtrl)
	circleCtrl.ConnectMousePressEvent(evh.HandleMousePressEvent)
	circleCtrl.ConnectMouseReleaseEvent(evh.HandleMouseReleaseEvent)
	pg.sceneItems.SetControlCircle(knotNo, isEntry, circleCtrl)
	vtx, vty := vertex.Coord()
	pg.sceneItems.SetControlLine(knotNo, isEntry, vtx, vty, ctrlx, ctrly, pg.penCtrl)
}

func (pg *Playground) addSegmentPaths(fromSegmentNo int, toSegmentNo int, pen gui.QPen_ITF) {
	paco := NewQPathCollector2d()
	pg.spline.Approx(fromSegmentNo, toSegmentNo, 0.5, paco)
	fmt.Printf("#line-segments: %v \n", paco.LineCnt())
	for segmNo := fromSegmentNo; segmNo <= toSegmentNo; segmNo++ {
		pg.sceneItems.SetSegmentPath(segmNo, paco.Paths[segmNo], pen, gui.NewQBrush())
	}
}

func (pg *Playground) addSplineToScene() {
	// bezier-control as solid gray circle

	// vertices
	knots := pg.spline.Knots()
	for i := 0; i < knots.KnotCnt(); i++ {
		t, _ := knots.Knot(i)

		x, y := pg.spline.At(t)
		pg.addVertexToScene(i, x, y)

		// controls
		switch spl := pg.spline.(type) {
		case *cubic.BezierSpline2d:
			// bezier control points
			bvt, _ := spl.Vertex(i).(*cubic.BezierVx2)
			entry := bvt.Entry()
			pg.addControlPointToScene(i, bvt, entry.X(), entry.Y(), true)
			exit := bvt.Exit()
			pg.addControlPointToScene(i, bvt, exit.X(), exit.Y(), false)
		case *cubic.HermiteSpline2d, *cubic.NaturalHermiteSpline2d:
			hvt, _ := spl.Vertex(i).(*cubic.HermiteVx2)
			entry := hvt.Control(true)
			pg.addControlPointToScene(i, hvt, entry.X(), entry.Y(), true)
			exit := hvt.Control(false)
			pg.addControlPointToScene(i, hvt, exit.X(), exit.Y(), false)
		default:
			panic(fmt.Sprintf("type not yet supported: %T", spl))
		}
	}

	// line segments
	pg.addSegmentPaths(0, pg.spline.Knots().SegmentCnt()-1, pg.pen)
}

type VertexEventHandler struct {
	playground *Playground
	knotNo     int
}

func (eh *VertexEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	//eh.mousePressX, eh.mousePressY = event.Pos().X(), event.Pos().Y()
	//fmt.Printf("mouse-press-event for vertex with knotNo = %v at %v/%v\n", eh.knotNo, eh.mousePressX, eh.mousePressY)
}

func (eh *VertexEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	x, y := pos.X(), pos.Y()
	vt := eh.playground.spline.Vertex(eh.knotNo)
	xold, yold := vt.Coord()
	/*fmt.Printf("mouse-released-event for vertex with knotNo = %v at %v/%v, for knot previously at %v/%v\n",
	eh.knotNo, x, y, xold, yold)*/

	// modify spline
	vt = vt.Translate(x-xold, y-yold)
	eh.playground.spline.UpdateVertex(eh.knotNo, vt)

	// move vertex
	circleVx := eh.playground.sceneItems.VertexCircle(eh.knotNo)
	circleVx.SetRect(eh.playground.vertexRectForCircle(x, y))

	// move control-points
	moveControlPoint := func(isEntry bool) {
		var ctrl *cubic.Control
		switch vertex := vt.(type) {
		case *cubic.BezierVx2:
			ctrl = vertex.Control(isEntry)
		case *cubic.HermiteVx2:
			ctrl = vertex.Control(isEntry)
		}
		circleEntry := eh.playground.sceneItems.ControlCircle(eh.knotNo, isEntry)
		if circleEntry != nil {
			circleEntry.SetRect(eh.playground.controlRectForCircle(ctrl.X(), ctrl.Y()))
			eh.playground.sceneItems.SetControlLine(eh.knotNo, isEntry, x, y, ctrl.X(), ctrl.Y(), eh.playground.penCtrl)
		}
	}
	moveControlPoint(true)
	moveControlPoint(false)

	// redraw segment paths
	fromSegmentNo, toSegmentNo, _ := bendit.AdjacentSegments(eh.playground.spline.Knots(), eh.knotNo, true, true)
	eh.playground.addSegmentPaths(fromSegmentNo, toSegmentNo, gui.NewQPen3(gui.NewQColor2(core.Qt__black)))
}

// incomplete: for simplicity reasons only add by double-click on the last vertex and delete by double-click on the first
func (eh *VertexEventHandler) HandleMouseDoubleClickEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	x, y := pos.X(), pos.Y()
	fmt.Printf("mouse-double-click-event for vertex with knotNo = %v at %v/%v, for knot previously at %v/%v\n",
		eh.knotNo, x, y)

	if eh.knotNo == eh.playground.spline.Knots().KnotCnt()-1 {
		// double-click on last vertex => add new one
		// TODO support HermiteVx2 also
		newBezierVx := cubic.NewBezierVx2(x+30, y+30, cubic.NewControl(x-20, y-20), nil)
		newKnotNo := eh.knotNo + 1
		eh.playground.spline.AddVertex(newKnotNo, newBezierVx)
		x, y = newBezierVx.Coord()
		eh.playground.addVertexToScene(newKnotNo, x, y)
		eh.playground.addControlPointToScene(newKnotNo, newBezierVx, newBezierVx.Entry().X(), newBezierVx.Entry().Y(), true)
		eh.playground.addControlPointToScene(newKnotNo, newBezierVx, newBezierVx.Entry().X(), newBezierVx.Entry().Y(), false)
	}
}

type ControlPointEventHandler struct {
	playground *Playground
	knotNo     int
	isEntry    bool
}

func (eh *ControlPointEventHandler) HandleMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
}

func (eh *ControlPointEventHandler) HandleMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	pos := event.Pos()
	x, y := pos.X(), pos.Y()
	vt := eh.playground.spline.Vertex(eh.knotNo)
	vtx, vty := vt.Coord()
	var dependent bool
	var othx, othy float64

	// modify spline
	switch vertex := vt.(type) {
	case *cubic.BezierVx2:
		ctrl := cubic.NewControl(x, y)
		if eh.isEntry {
			vt = vertex.WithEntry(ctrl)
		} else {
			vt = vertex.WithExit(ctrl)
		}
		dependent = vertex.Dependent()
		if dependent {
			otherCtrl := vt.(*cubic.BezierVx2).Control(!eh.isEntry)
			othx, othy = otherCtrl.X(), otherCtrl.Y()
		}
	case *cubic.HermiteVx2:
		if eh.isEntry {
			vt = vertex.WithEntryTan(cubic.NewControl(vtx-x, vty-y))
		} else {
			vt = vertex.WithExitTan(cubic.NewControl(x-vtx, y-vty))
		}
		dependent = vertex.Dependent()
		if dependent {
			otherCtrl := vt.(*cubic.HermiteVx2).Control(!eh.isEntry)
			othx, othy = otherCtrl.X(), otherCtrl.Y()
		}
	}
	eh.playground.spline.UpdateVertex(eh.knotNo, vt)

	// move control
	ctrlCircle := eh.playground.sceneItems.ControlCircle(eh.knotNo, eh.isEntry)
	ctrlCircle.SetRect(eh.playground.controlRectForCircle(x, y))
	eh.playground.sceneItems.SetControlLine(eh.knotNo, eh.isEntry, vtx, vty, x, y, eh.playground.penCtrl)

	if dependent {
		ctrlCircle = eh.playground.sceneItems.ControlCircle(eh.knotNo, !eh.isEntry)
		if ctrlCircle != nil {
			ctrlCircle.SetRect(eh.playground.controlRectForCircle(othx, othy))
			eh.playground.sceneItems.SetControlLine(eh.knotNo, !eh.isEntry, vtx, vty, othx, othy, eh.playground.penCtrl)
		}
	}

	// replace segment paths (on both side of vertex if dependent)
	fromSegmentNo, toSegmentNo, _ := bendit.AdjacentSegments(eh.playground.spline.Knots(), eh.knotNo,
		eh.isEntry || dependent, !eh.isEntry || dependent)
	eh.playground.addSegmentPaths(fromSegmentNo, toSegmentNo, gui.NewQPen3(gui.NewQColor2(core.Qt__black)))
}

/*func (pg *Playground) drawSplineBySubdivisionPath(qp *gui.QPainter) {
	paco := NewQPathCollector2d()
	bendit.ApproxAll(pg.spline, 0.5, paco)
	fmt.Printf("#line-segments: %v \n", paco.LineCnt())
	pen := gui.NewQPen()
	for _, path := range paco.Paths {
		qp.StrokePath(path, pen)
	}
}*/

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
