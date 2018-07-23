package cairo

type filler struct {
	polygon      *polygon
	tolerance    float64
	limit        box
	hasLimits    bool
	currentPoint point
	lastMoveTo   point
}

func fillerLineTo(closure interface{}, point *point) Status {
	filler := closure.(*filler)
	status := filler.polygon.addExternalEdge(&filler.currentPoint, point)
	filler.currentPoint = *point
	return status
}

func fillerClose(closure interface{}) Status {
	filler := closure.(*filler)
	/* close the subpath */
	return fillerLineTo(closure, &filler.lastMoveTo)
}

func fillerMoveTo(closure interface{}, point *point) Status {
	filler := closure.(*filler)

	/* close the subpath */
	status := fillerClose(closure)
	if status != 0 {
		return status
	}

	/* make sure that the closure represents a degenerate path */
	filler.currentPoint = *point
	filler.lastMoveTo = *point
	return StatusSuccess
}

func fillerCurveTo(closure interface{}, p1, p2, p3 *point) Status {
	filler := closure.(*filler)

	if filler.hasLimits {
		if !splineIntersects(&filler.currentPoint, p1, p2, p3, &filler.limit) {
			return fillerLineTo(closure, p3)
		}
	}

	addPointFunc := func(closure interface{}, point *point, tangent *slope) Status {
		return fillerLineTo(closure, point)
	}
	var spline spline
	if !spline.init(addPointFunc, filler, &filler.currentPoint, p1, p2, p3) {
		return fillerLineTo(closure, p3)
	}

	return spline.decompose(filler.tolerance)
}

func (path *pathFixed) fillToPolygon(tolerance float64, polygon *polygon) Status {
	var filler filler
	filler.polygon = polygon
	filler.tolerance = tolerance

	filler.hasLimits = false
	if len(polygon.limits) != 0 {
		filler.hasLimits = true
		filler.limit = polygon.limit
	}

	/* make sure that the closure represents a degenerate path */
	filler.currentPoint.x = 0
	filler.currentPoint.y = 0
	filler.lastMoveTo = filler.currentPoint

	status := path.interpret(fillerMoveTo, fillerLineTo, fillerCurveTo, fillerClose, &filler)
	if status != 0 {
		return status
	}

	return fillerClose(&filler)
}

type fillerRectilinearAligned struct {
	polygon      *polygon
	currentPoint point
	lastMoveTo   point
}

func fillerRALineTo(closure interface{}, point0 *point) Status {
	filler := closure.(*fillerRectilinearAligned)
	var p point
	p.x = point0.x.roundDown()
	p.y = point0.y.roundDown()

	status := filler.polygon.addExternalEdge(&filler.currentPoint, &p)
	filler.currentPoint = p
	return status
}

func fillerRAClose(closure interface{}) Status {
	filler := closure.(*fillerRectilinearAligned)
	return fillerRALineTo(closure, &filler.lastMoveTo)
}

func fillerRAMoveTo(closure interface{}, point0 *point) Status {
	filler := closure.(*fillerRectilinearAligned)
	var p point

	/* close current subpath */
	status := fillerRAClose(closure)
	if status != 0 {
		return status
	}

	p.x = point0.x.roundDown()
	p.y = point0.y.roundDown()

	/* make sure that the closure represents a degenerate path */
	filler.currentPoint = p
	filler.lastMoveTo = p
	return StatusSuccess
}

func (path *pathFixed) fillRectilinearToPolygon(antialias Antialias, polygon *polygon) Status {
	var filler fillerRectilinearAligned
	if antialias != AntialiasNone {
		return path.fillToPolygon(0, polygon)
	}

	filler.polygon = polygon

	/* make sure that the closure represents a degenerate path */
	filler.currentPoint.x = 0
	filler.currentPoint.y = 0
	filler.lastMoveTo = filler.currentPoint

	status := path.interpretFlat(fillerRAMoveTo, fillerRALineTo, fillerRAClose, &filler, 0)
	if status != 0 {
		return status
	}
	return fillerRAClose(&filler)
}

func (path *pathFixed) fillToTraps(fillRule FillRule, tolerance float64, traps *traps) Status {
	var polygon polygon
	if path.fillIsEmpty() {
		return StatusSuccess
	}

	polygon.init(traps.limits)
	status := path.fillToPolygon(tolerance, &polygon)
	if status != 0 || len(polygon.edges) == 0 {
		goto CLEANUP
	}

	status = bentleyOttmannTessellatePolygon(traps, &polygon, fillRule)

CLEANUP:
	polygon.fini()
	return status
}

func (path *pathFixed) fillRectilinearTessellateToBoxes(fillRule FillRule, antialias Antialias, boxes *boxes) Status {
	var polygon polygon
	polygon.init(boxes.limits)
	boxes.limits = nil

	/* tolerance will be ignored as the path is rectilinear */
	status := path.fillRectilinearToPolygon(antialias, &polygon)
	if status == StatusSuccess {
		status = bentleyOttmannTessellateRectilinearPolygonToBoxes(&polygon, fillRule, boxes)
	}
	polygon.fini()
	return status
}

func (path *pathFixed) fillRectilinearToBoxes(fillRule FillRule, antialias Antialias, boxes *boxes) Status {
	var iter pathFixedIter
	var box box
	var status Status

	if path.isBox(&box) {
		return boxes.add(antialias, &box)
	}

	iter.init(path)
	for iter.isFillBox(&box) {
		if box.p1.y == box.p2.y || box.p1.x == box.p2.x {
			continue
		}

		if box.p1.y > box.p2.y {
			t := box.p1.y
			box.p1.y = box.p2.y
			box.p2.y = t

			t = box.p1.x
			box.p1.x = box.p2.x
			box.p2.x = t
		}

		status = boxes.add(antialias, &box)
		if status != 0 {
			return status
		}
	}

	if iter.addEnd() {
		return bentleyOttmannTessellateBoxes(boxes, fillRule, boxes)
	}

	/* path is not rectangular, try extracting clipped rectilinear edges */
	boxes.clear()
	return path.fillRectilinearTessellateToBoxes(fillRule, antialias, boxes)
}
