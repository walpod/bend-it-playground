package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/walpod/bendigo"
	"github.com/walpod/bendigo/cubic"
)

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

func (si *GraphicsSceneItems) SetControlLine(knotNo int, isEntry bool, from, to bendigo.Vec, pen gui.QPen_ITF) *widgets.QGraphicsLineItem {
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
	lineItem := si.scene.AddLine2(from[0], from[1], to[0], to[1], pen)
	si.controlLines[knotNo][sideNo] = lineItem
	return lineItem
}

func (si *GraphicsSceneItems) Clear() {
	si.scene.Clear()
	si.segmentPaths = nil
	si.vertexCircles = nil
	si.controlCircles = nil
	si.controlLines = nil
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
	splineBuilder bendigo.SplineVertBuilder
	sceneItems    GraphicsSceneItems
	// styles for spline and vertices
	pen   gui.QPen_ITF
	brush gui.QBrush_ITF
	// styles for controls
	penCtrl           gui.QPen_ITF
	brushCtrl         gui.QBrush_ITF
	brushCtrlLeading  gui.QBrush_ITF
	brushCtrlDisabled gui.QBrush_ITF
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
	//darkGray := gui.NewQColor2(core.Qt__darkGray)
	white := gui.NewQColor2(core.Qt__white)
	pg.pen = gui.NewQPen3(black)
	pg.brush = gui.NewQBrush2(core.Qt__SolidPattern)
	pg.penCtrl = gui.NewQPen2(core.Qt__DotLine) // core.Qt__DashLine
	pg.brushCtrl = gui.NewQBrush3(white, core.Qt__SolidPattern)
	pg.brushCtrlLeading = gui.NewQBrush3(gray, core.Qt__SolidPattern)
	pg.brushCtrlDisabled = gui.NewQBrush3(white, core.Qt__SolidPattern)

	pg.prepareSplineBuilder()
	pg.addSplineToScene()
	return pg
}

func (pg *Playground) prepareSplineBuilder() {
	// hermite
	/*pg.splineBuilder = cubic.NewHermiteVertBuilder(nil,
		cubic.NewHermiteVertex(bendigo.NewVec(200, 200), nil, bendigo.NewVec(90, 90)),
		cubic.NewHermiteVertex(bendigo.NewVec(350, 350), bendigo.NewVec(200, 0), nil),
		cubic.NewHermiteVertex(bendigo.NewVec(500, 200), bendigo.NewVec(100, -100), nil),
	)*/

	/*pg.splineBuilder = cubic.NewNaturalVertBuilder(nil,
		cubic.NewRawHermiteVertex(bendigo.NewVec(10, 10)),
		cubic.NewRawHermiteVertex(bendigo.NewVec(100, 100)),
		cubic.NewRawHermiteVertex(bendigo.NewVec(150, 10)),
	)*/

	/*nat := cubic.NewNaturalVertBuilder(nil)
	nat.AddVertex(0, cubic.NewRawHermiteVertex(bendigo.NewVec(100, 100)))
	nat.AddVertex(1, cubic.NewRawHermiteVertex(bendigo.NewVec(400, 400)))
	nat.AddVertex(2, cubic.NewRawHermiteVertex(bendigo.NewVec(700, 100)))
	pg.splineBuilder = nat*/

	/*pg.splineBuilder = cubic.NewCardinalVertBuilder(nil, -3,
		cubic.NewRawHermiteVertex(bendigo.NewVec(100, 100)),
		cubic.NewRawHermiteVertex(bendigo.NewVec(400, 400)),
		cubic.NewRawHermiteVertex(bendigo.NewVec(700, 100)),
	)*/

	// bezier
	/*pg.splineBuilder = cubic.NewBezierVertBuilder(nil,
	cubic.NewBezierVertex(bendigo.NewVec(200, 200), nil, bendigo.NewVec(250, 200)),
	cubic.NewBezierVertex(bendigo.NewVec(400, 400), bendigo.NewVec(350, 400), nil))*/

	/*pg.splineBuilder = cubic.NewBezierVertBuilder(nil,
	cubic.NewBezierVertex(bendigo.NewVec(200, 200), bendigo.NewVec(100, 200), bendigo.NewVec(300, 200)),
	cubic.NewBezierVertex(bendigo.NewVec(300, 300), bendigo.NewVec(200, 300), bendigo.NewVec(400, 300)))*/

	pg.splineBuilder = cubic.NewBezierVertBuilder(nil)
	pg.splineBuilder.AddVertex(0, cubic.NewBezierVertex(bendigo.NewVec(100, 100), nil, bendigo.NewVec(120, 150)))
	pg.splineBuilder.AddVertex(1, cubic.NewBezierVertex(bendigo.NewVec(300, 300), bendigo.NewVec(200, 300), nil))
	pg.splineBuilder.AddVertex(2, cubic.NewBezierVertex(bendigo.NewVec(500, 100), bendigo.NewVec(490, 150), nil))
}

func (pg *Playground) HasAutoControls() bool {
	switch pg.splineBuilder.(type) {
	case *cubic.NaturalVertBuilder, *cubic.CardinalVertBuilder:
		// on change of a vertex natural and cardinal splines are changing controls of other vertices on-the-fly
		return true
	default:
		return false
	}
}

func (pg *Playground) vertexRectForCircle(x float64, y float64) *core.QRectF {
	radius := 6.0
	return core.NewQRectF4(x-radius, y-radius, 2*radius, 2*radius)
}

func (pg *Playground) addVertexToScene(knotNo int, v bendigo.Vec) {
	veh := VertexEventHandler{playground: pg, knotNo: knotNo}
	// vertex as solid black circle
	circleVt := widgets.NewQGraphicsEllipseItem2(pg.vertexRectForCircle(v[0], v[1]), nil)
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

func (pg *Playground) addControlPointToScene(knotNo int, vertex bendigo.Vertex, ctrl bendigo.Vec, isEntry bool) {
	evh := ControlPointEventHandler{playground: pg, knotNo: knotNo, isEntry: isEntry}
	// control as solid gray circle
	circleCtrl := widgets.NewQGraphicsEllipseItem2(pg.controlRectForCircle(ctrl[0], ctrl[1]), nil)
	if pg.HasAutoControls() {
		circleCtrl.SetBrush(pg.brushCtrlDisabled)
	} else {
		ev := vertex.(*cubic.EnexVertex)
		if ev.Leading() {
			circleCtrl.SetBrush(pg.brushCtrlLeading)
		} else {
			circleCtrl.SetBrush(pg.brushCtrl)
		}
	}
	circleCtrl.ConnectMousePressEvent(evh.HandleMousePressEvent)
	circleCtrl.ConnectMouseReleaseEvent(evh.HandleMouseReleaseEvent)
	circleCtrl.ConnectMouseDoubleClickEvent(evh.HandleMouseDoubleClickEvent)
	pg.sceneItems.SetControlCircle(knotNo, isEntry, circleCtrl)
	pg.sceneItems.SetControlLine(knotNo, isEntry, vertex.Loc(), ctrl, pg.penCtrl)
}

func (pg *Playground) addSegmentPaths(fromSegmentNo int, toSegmentNo int, pen gui.QPen_ITF) {
	paco := NewQPathCollector()
	pg.splineBuilder.LinApproximate(fromSegmentNo, toSegmentNo, paco, bendigo.NewLinaxParams(0.5))
	fmt.Printf("#line-segments: %v \n", paco.LineCnt())
	for segmNo := fromSegmentNo; segmNo <= toSegmentNo; segmNo++ {
		pg.sceneItems.SetSegmentPath(segmNo, paco.Paths[segmNo], pen, gui.NewQBrush())
	}
}

func (pg *Playground) addSplineToScene() {
	// bezier-control as solid gray circle

	// vertices
	vertices := bendigo.Vertices(pg.splineBuilder)
	for i := 0; i < len(vertices); i++ {
		pg.addVertexToScene(i, vertices[i].Loc())

		// controls
		vt, _ := pg.splineBuilder.Vertex(i).(*cubic.EnexVertex)
		pg.addControlPointToScene(i, vt, vt.ControlAsAbsolute(true), true)
		pg.addControlPointToScene(i, vt, vt.ControlAsAbsolute(false), false)
	}

	// line segments
	pg.addSegmentPaths(0, pg.splineBuilder.Knots().SegmentCnt()-1, pg.pen)
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
	loc := bendigo.NewVec(pos.X(), pos.Y())
	vt := eh.playground.splineBuilder.Vertex(eh.knotNo).(*cubic.EnexVertex)
	oldLoc := vt.Loc()
	/*fmt.Printf("mouse-released-event for vertex with knotNo = %v at %v/%v, for knot previously at %v/%v\n",
	eh.knotNo, x, y, xold, yold)*/

	// modify spline
	vt.Shift(loc.Sub(oldLoc))
	eh.playground.splineBuilder.UpdateVertex(eh.knotNo, vt)

	if eh.playground.HasAutoControls() {
		eh.playground.sceneItems.Clear()
		eh.playground.addSplineToScene()
		return
	}

	// shift vertex-circle
	circleVx := eh.playground.sceneItems.VertexCircle(eh.knotNo)
	circleVx.SetRect(eh.playground.vertexRectForCircle(loc[0], loc[1]))

	// shift control-circles
	shiftControlCircle := func(isEntry bool) {
		ctrlLoc := vt.ControlAsAbsolute(isEntry)
		circleEntry := eh.playground.sceneItems.ControlCircle(eh.knotNo, isEntry)
		if circleEntry != nil {
			circleEntry.SetRect(eh.playground.controlRectForCircle(ctrlLoc[0], ctrlLoc[1]))
			eh.playground.sceneItems.SetControlLine(eh.knotNo, isEntry, loc, ctrlLoc, eh.playground.penCtrl)
		}
	}
	shiftControlCircle(true)
	shiftControlCircle(false)

	// redraw segment paths
	fromSegmentNo, toSegmentNo, _ := bendigo.SegmentsAroundKnot(eh.playground.splineBuilder.Knots(), eh.knotNo, true, true)
	eh.playground.addSegmentPaths(fromSegmentNo, toSegmentNo, gui.NewQPen3(gui.NewQColor2(core.Qt__black)))
}

// incomplete: for simplicity reasons only add by double-click on the last vertex and delete by double-click on the first
func (eh *VertexEventHandler) HandleMouseDoubleClickEvent(event *widgets.QGraphicsSceneMouseEvent) {
	/*pos := event.Pos()
	x, y := pos.X(), pos.Y()
	fmt.Printf("mouse-double-click-event for vertex with knotNo = %v at %v/%v, for knot previously at %v/%v\n",
		eh.knotNo, x, y)*/

	if eh.knotNo == eh.playground.splineBuilder.Knots().KnotCnt()-1 {
		// double-click on last vertex => add new one
		vt := eh.playground.splineBuilder.Vertex(eh.knotNo).(*cubic.EnexVertex)
		vt.Shift(bendigo.NewVec(30, 30))
		newKnotNo := eh.knotNo + 1
		eh.playground.splineBuilder.AddVertex(newKnotNo, vt)

		// draw
		if eh.playground.HasAutoControls() {
			eh.playground.sceneItems.Clear()
			eh.playground.addSplineToScene()
			return
		}
		eh.playground.addVertexToScene(newKnotNo, vt.Loc())
		eh.playground.addControlPointToScene(newKnotNo, vt, vt.ControlAsAbsolute(true), true)
		eh.playground.addControlPointToScene(newKnotNo, vt, vt.ControlAsAbsolute(false), false)
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
	ctrlLoc := bendigo.NewVec(pos.X(), pos.Y())
	vt := eh.playground.splineBuilder.Vertex(eh.knotNo).(*cubic.EnexVertex)

	// don't allow to move controls directly if they are calculated
	if eh.playground.HasAutoControls() {
		return
	}

	// modify spline
	vt.SetControl(ctrlLoc, eh.isEntry)
	eh.playground.splineBuilder.UpdateVertex(eh.knotNo, vt)

	// shift control-circle
	ctrlCircle := eh.playground.sceneItems.ControlCircle(eh.knotNo, eh.isEntry)
	ctrlCircle.SetRect(eh.playground.controlRectForCircle(ctrlLoc[0], ctrlLoc[1]))
	eh.playground.sceneItems.SetControlLine(eh.knotNo, eh.isEntry, vt.Loc(), ctrlLoc, eh.playground.penCtrl)

	if vt.Leading() {
		ctrlCircle = eh.playground.sceneItems.ControlCircle(eh.knotNo, !eh.isEntry)
		if ctrlCircle != nil {
			otherCtrlLoc := vt.ControlAsAbsolute(!eh.isEntry)
			ctrlCircle.SetRect(eh.playground.controlRectForCircle(otherCtrlLoc[0], otherCtrlLoc[1]))
			eh.playground.sceneItems.SetControlLine(eh.knotNo, !eh.isEntry, vt.Loc(), otherCtrlLoc, eh.playground.penCtrl)
		}
	}

	// replace segment paths (on both side of vertex if 'leading')
	fromSegmentNo, toSegmentNo, _ := bendigo.SegmentsAroundKnot(eh.playground.splineBuilder.Knots(), eh.knotNo,
		eh.isEntry || vt.Leading(), !eh.isEntry || vt.Leading())
	eh.playground.addSegmentPaths(fromSegmentNo, toSegmentNo, gui.NewQPen3(gui.NewQColor2(core.Qt__black)))
}

func (eh ControlPointEventHandler) HandleMouseDoubleClickEvent(event *widgets.QGraphicsSceneMouseEvent) {
	// toggle 'Leading' property of vertex
	vt := eh.playground.splineBuilder.Vertex(eh.knotNo).(*cubic.EnexVertex)
	vt.ToggleLeading(eh.isEntry)

	// change color of control-circles
	eh.playground.addControlPointToScene(eh.knotNo, vt, vt.ControlAsAbsolute(true), true)
	eh.playground.addControlPointToScene(eh.knotNo, vt, vt.ControlAsAbsolute(false), false)

	// if changed to 'Leading' then replace segment path on follower side
	if vt.Leading() {
		fromSegmentNo, toSegmentNo, _ := bendigo.SegmentsAroundKnot(eh.playground.splineBuilder.Knots(), eh.knotNo,
			vt.EntryLeads(), !vt.EntryLeads())
		eh.playground.addSegmentPaths(fromSegmentNo, toSegmentNo, gui.NewQPen3(gui.NewQColor2(core.Qt__black)))
	}
}

type QPathCollector struct {
	Paths map[int]*gui.QPainterPath
}

func NewQPathCollector() *QPathCollector {
	return &QPathCollector{Paths: map[int]*gui.QPainterPath{}}
}

//segmentNo int, tstart, tend, pstartx, pstarty, pendx, pendy float64
func (lc *QPathCollector) ConsumeLine(segmentNo int, tstart, tend float64, pstart, pend bendigo.Vec) {
	// get path for segment
	path, exists := lc.Paths[segmentNo]
	if !exists {
		path = gui.NewQPainterPath()
		lc.Paths[segmentNo] = path
	}

	// add line to path
	if path.ElementCount() == 0 {
		path.MoveTo(core.NewQPointF3(pstart[0], pstart[1]))
	}
	path.LineTo(core.NewQPointF3(pend[0], pend[1]))
}

func (lc *QPathCollector) LineCnt() int {
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
