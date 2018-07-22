package cairo

func (path *pathFixed) init() {
	path.currentPoint.x = 0
	path.currentPoint.y = 0
	path.lastMovePoint = path.currentPoint

	path.hasCurrentPoint = false
	path.needsMoveTo = true
	path.hasExtents = false
	path.hasCurveTo = false
	path.strokeIsRectilinear0 = true
	path.fillIsRectilinear0 = true
	path.fillMaybeRegion0 = true
	path.fillIsEmpty0 = true

	path.extents.p1.x = 0
	path.extents.p1.y = 0
	path.extents.p1.x = 0
	path.extents.p2.y = 0
}

func (path *pathFixed) initCopy(other *pathFixed) Status {
	path.currentPoint = other.currentPoint
	path.lastMovePoint = other.lastMovePoint

	path.hasCurrentPoint = other.hasCurrentPoint
	path.needsMoveTo = other.needsMoveTo
	path.hasExtents = other.hasExtents
	path.hasCurveTo = other.hasCurveTo
	path.strokeIsRectilinear0 = other.strokeIsRectilinear0
	path.fillIsRectilinear0 = other.fillIsRectilinear0
	path.fillMaybeRegion0 = other.fillMaybeRegion0
	path.fillIsEmpty0 = other.fillIsEmpty0

	path.extents = other.extents

	path.ops = make([]pathOp, len(other.ops))
	copy(path.ops, other.ops)
	path.points = make([]point, len(other.points))
	copy(path.points, other.points)
	return StatusSuccess
}

func (a *pathFixed) equal(b *pathFixed) bool {
	if a == b {
		return true
	}

	/* use the flags to quickly differentiate based on contents */
	if a.hasCurveTo != b.hasCurveTo {
		return false
	}

	if a.extents.p1.x != b.extents.p1.x ||
		a.extents.p1.y != b.extents.p1.y ||
		a.extents.p2.x != b.extents.p2.x ||
		a.extents.p2.y != b.extents.p2.y {
		return false
	}

	if len(a.ops) == 0 && len(b.ops) == 0 {
		return true
	}

	if len(a.ops) != len(b.ops) || len(a.points) != len(b.points) {
		return false
	}

	for i := 0; i < len(a.ops); i++ {
		if a.ops[i] != b.ops[i] {
			return false
		}
	}

	for i := 0; i < len(a.points); i++ {
		if a.points[i] != b.points[i] {
			return false
		}
	}

	return true
}

func pathFixedCreate() *pathFixed {
	path := new(pathFixed)
	path.init()
	return path
}

func (path *pathFixed) fini() {
	return
}

func (path *pathFixed) destroy() {
	return
}

func (path *pathFixed) lastOp() pathOp {
	return path.ops[len(path.ops)-1]
}

func (path *pathFixed) penultimatePoint() *point {
	return &path.points[len(path.points)-2]
}

func (path *pathFixed) dropLineTo() {
	if path.lastOp() != pathOpLineTo {
		panic("assert failed path.lastOp == pathOpLineTo")
	}

	path.ops = path.ops[:len(path.ops)-1]
	path.points = path.points[:len(path.points)-1]
}

func (path *pathFixed) moveTo(x, y fixed) Status {
	path.newSubPath()

	path.hasCurrentPoint = true
	path.currentPoint.x = x
	path.currentPoint.y = y
	path.lastMovePoint = path.currentPoint
	return StatusSuccess
}

func (path *pathFixed) moveToApply() Status {
	if !path.needsMoveTo {
		return StatusSuccess
	}

	path.needsMoveTo = false

	if path.hasExtents {
		path.extents.addPoint(&path.currentPoint)
	} else {
		path.extents.set(&path.currentPoint, &path.currentPoint)
		path.hasExtents = true
	}

	if path.fillMaybeRegion0 {
		path.fillMaybeRegion0 = path.currentPoint.x.IsInteger() &&
			path.currentPoint.y.IsInteger()
	}
	path.lastMovePoint = path.currentPoint
	return path.add(pathOpMoveTo, []point{path.currentPoint})
}

func (path *pathFixed) newSubPath() {
	if !path.needsMoveTo {
		/* If the current subpath doesn't need_move_to, it contains at least one command */
		if path.fillIsRectilinear0 {
			/* Implicitly close for fill */
			path.fillIsRectilinear0 = path.currentPoint.x == path.lastMovePoint.x ||
				path.currentPoint.y == path.lastMovePoint.y

			path.fillMaybeRegion0 = path.fillMaybeRegion0 && path.fillIsRectilinear0
		}
		path.needsMoveTo = true
	}
	path.hasCurrentPoint = false
}

func (path *pathFixed) relMoveTo(dx, dy fixed) Status {
	if !path.hasCurrentPoint {
		return StatusNoCurrentPoint.error()
	}

	return path.moveTo(path.currentPoint.x+dx,
		path.currentPoint.y+dy)
}

func (path *pathFixed) lineTo(x, y fixed) Status {
	var status Status
	var point0 point

	point0.x = x
	point0.y = y

	/* When there is not yet a current point, the line_to operation
	 * becomes a move_to instead. Note: We have to do this by
	 * explicitly calling into _cairo_path_fixed_move_to to ensure
	 * that the last_move_point state is updated properly.
	 */
	if !path.hasCurrentPoint {
		return path.moveTo(point0.x, point0.y)
	}
	status = path.moveToApply()
	if status != 0 {
		return status
	}

	/* If the previous op was but the initial MOVE_TO and this segment
	 * is degenerate, then we can simply skip this point. Note that
	 * a move-to followed by a degenerate line-to is a valid path for
	 * stroking, but at all other times is simply a degenerate segment.
	 */
	if path.lastOp() != pathOpMoveTo {
		if x == path.currentPoint.x && y == path.currentPoint.y {
			return StatusSuccess
		}
	}

	/* If the previous op was also a LINE_TO with the same gradient,
	 * then just change its end-point rather than adding a new op.
	 */
	if path.lastOp() == pathOpLineTo {
		var p *point
		p = path.penultimatePoint()
		if p.x == path.currentPoint.x && p.y == path.currentPoint.y {
			/* previous line element was degenerate, replace */
			path.dropLineTo()
		} else {
			var prev, self slope
			prev.init(p, &path.currentPoint)
			self.init(&path.currentPoint, &point0)

			if prev.equal(&self) &&
				/* cannot trim anti-parallel segments whilst stroking */
				!prev.backwards(&self) {
				path.dropLineTo()
				/* In this case the flags might be more restrictive than
				 * what we actually need.
				 * When changing the flags definition we should check if
				 * changing the line_to point can affect them.
				 */
			}
		}
	}

	if path.strokeIsRectilinear0 {
		path.strokeIsRectilinear0 = path.currentPoint.x == x ||
			path.currentPoint.y == y
		path.fillIsRectilinear0 = path.fillIsRectilinear0 && path.strokeIsRectilinear0
		path.fillMaybeRegion0 = path.fillMaybeRegion0 && path.fillIsRectilinear0
		if path.fillMaybeRegion0 {
			path.fillMaybeRegion0 = x.IsInteger() && y.IsInteger()
		}
		if path.fillIsEmpty0 {
			path.fillIsEmpty0 = path.currentPoint.x == x &&
				path.currentPoint.y == y
		}
	}

	path.currentPoint = point0
	path.extents.addPoint(&point0)
	return path.add(pathOpLineTo, []point{point0})
}

var _ = `
`

func (path *pathFixed) isBox(box *box) bool {
	var buf *pathBuf
	if !path.fillIsRectilinear0 {
		return false
	}

	if !path.isQuad() {
		return false
	}

	buf = path.head()
	if pointsFromRect(buf.points) {
		canonicalBox(box, &buf.points[0], &buf.points[2])
		return true
	}

	return false
}
