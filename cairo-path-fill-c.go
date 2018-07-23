package cairo

type filler struct {
	polygon      *polygon
	tolearance   float64
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

	return spline.decompose(filler.tolearance)
}

func (p *pathFixed) fillRectilinearToBoxes(fillRule FillRule, antialias Antialias, boxes *boxes) Status {
	// TODO
	return StatusSuccess
}
