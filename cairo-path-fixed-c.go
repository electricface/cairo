package cairo

func (path *pathFixed) init() {
	path.buf = []*pathBuf{
		{
			ops:    path.buf0.ops[:0],
			points: path.buf0.points[:0],
		},
	}

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

	path.buf = make([]*pathBuf, len(other.buf))
	path.buf[0] = &pathBuf{
		ops:    path.buf0.ops[:len(other.buf[0].ops)],
		points: path.buf0.points[:len(other.buf[0].points)],
	}
	for i := 1; i < len(other.buf); i++ {
		path.buf[i].ops = make([]pathOp, len(other.buf[i].ops))
		path.buf[i].points = make([]point, len(other.buf[i].points))
	}
	for i := 0; i < len(other.buf); i++ {
		copy(path.buf[i].ops, other.buf[i].ops)
		copy(path.buf[i].points, other.buf[i].points)
	}

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

	numOpsA := 0
	numPointsA := 0
	for _, buf := range a.buf {
		numOpsA += len(buf.ops)
		numPointsA += len(buf.points)
	}

	numOpsB := 0
	numPointsB := 0
	for _, buf := range b.buf {
		numOpsB += len(buf.ops)
		numPointsB += len(buf.points)
	}

	if numOpsA == 0 && numOpsB == 0 {
		return true
	}

	if numOpsA != numOpsB || numPointsA != numPointsB {
		return false
	}

	var bufAIdx, bufBIdx int
	bufA := a.buf[0]
	opsA := bufA.ops
	pointsA := bufA.points
	numOpsA = len(opsA)
	numPointsA = len(pointsA)

	bufB := b.buf[0]
	opsB := bufB.ops
	pointsB := bufB.points
	numOpsB = len(opsB)
	numPointsB = len(pointsB)

	for {
		numOps := minInt(numOpsA, numOpsB)
		numPoints := minInt(numPointsA, numPointsB)

		for i := 0; i < numOps; i++ {
			if opsA[i] != opsB[i] {
				return false
			}
		}

		for i := 0; i < numPoints; i++ {
			if pointsA[i] != pointsB[i] {
				return false
			}
		}

		numOpsA -= numOps
		opsA = opsA[:numOps]
		numPointsA -= numPoints
		pointsA = pointsA[:numPoints]

		if numOpsA == 0 || numPointsA == 0 {
			if numOpsA != 0 || numPointsA != 0 {
				return false
			}
			// next buf
			bufAIdx++
			if bufAIdx == len(a.buf) {
				break
			}
			bufA = a.buf[bufAIdx]

			pointsA = bufA.points
			numPointsA = len(pointsA)
			opsA = bufA.ops
			numOpsA = len(opsA)
		}

		numOpsB -= numOps
		opsB = opsB[:numOps]
		numPointsB -= numPoints
		pointsB = pointsB[:numPoints]

		if numOpsB == 0 || numPointsB == 0 {
			if numOpsB != 0 || numPointsB != 0 {
				return false
			}
			// next buf
			bufBIdx++
			if bufBIdx == len(b.buf) {
				break
			}
			bufB = b.buf[bufBIdx]

			pointsB = bufB.points
			numPointsB = len(pointsB)
			opsB = bufB.ops
			numOpsB = len(opsB)
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

func (path *pathFixed) head() *pathBuf {
	return path.buf[0]
}

func (path *pathFixed) tail() *pathBuf {
	return path.buf[len(path.buf)-1]
}

func (path *pathFixed) lastOp() pathOp {
	buf := path.tail()
	return buf.ops[len(buf.ops)-1]
}

func (path *pathFixed) penultimatePoint() *point {
	bufIdx := len(path.buf) - 1
	// buf is tail
	buf := path.buf[bufIdx]
	if len(buf.points) >= 2 {
		return &buf.points[len(buf.points)-2]
	} else {
		bufIdx--
		prevBuf := path.buf[bufIdx]
		if len(prevBuf.points) < 2-len(buf.points) {
			panic("assert failed len(prevBuf.points)	>= 2 - len(buf.points)")
		}
		return &prevBuf.points[len(prevBuf.points)-(2-len(buf.points))]
	}
}

func (path *pathFixed) dropLineTo() {
	if path.lastOp() != pathOpLineTo {
		panic("assert failed path.lastOp == pathOpLineTo")
	}

	buf := path.tail()
	buf.ops = buf.ops[:len(buf.ops)-1]
	buf.points = buf.points[:len(buf.points)-1]
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

func (path *pathFixed) relLineTo(dx, dy fixed) Status {
	if !path.hasCurrentPoint {
		return StatusNoCurrentPoint.error()
	}

	return path.lineTo(path.currentPoint.x+dx,
		path.currentPoint.y+dy)
}

func (path *pathFixed) curveTo(x0, y0, x1, y1, x2, y2 fixed) Status {
	var status Status
	var points [3]point

	/* If this curves does not move, replace it with a line-to.
	 * This frequently happens with rounded-rectangles and r==0.
	 */
	if path.currentPoint.x == x2 && path.currentPoint.y == y2 {
		if x1 == x2 && x0 == x2 && y1 == y2 && y0 == y2 {
			return path.lineTo(x2, y2)
		}
		/* We may want to check for the absence of a cusp, in which case
		 * we can also replace the curve-to with a line-to.
		 */
	}

	/* make sure subpaths are started properly */
	if !path.hasCurrentPoint {
		status = path.moveTo(x0, y0)
		if status != StatusSuccess {
			panic("assert failed status == StatusSuccess")
		}
	}

	status = path.moveToApply()
	if status != 0 {
		return status
	}

	/* If the previous op was a degenerate LINE_TO, drop it. */
	if path.lastOp() == pathOpLineTo {
		p := path.penultimatePoint()
		if p.x == path.currentPoint.x && p.y == path.currentPoint.y {
			/* previous line element was degenerate, replace */
			path.dropLineTo()
		}
	}

	points[0].x = x0
	points[0].y = y0

	points[1].x = x1
	points[1].y = y1

	points[2].x = x2
	points[2].y = y2

	path.extents.addCurveTo(&path.currentPoint, &points[0], &points[1], &points[2])

	path.currentPoint = points[2]
	path.hasCurveTo = true
	path.strokeIsRectilinear0 = false
	path.fillIsRectilinear0 = false
	path.fillMaybeRegion0 = false
	path.fillIsEmpty0 = false

	return path.add(pathOpCurveTo, points[:])
}

func (path *pathFixed) relCurveTo(dx0, dy0, dx1, dy1, dx2, dy2 fixed) Status {
	if !path.hasCurrentPoint {
		return StatusNoCurrentPoint.error()
	}

	return path.curveTo(path.currentPoint.x+dx0,
		path.currentPoint.y+dy0,
		path.currentPoint.x+dx1,
		path.currentPoint.y+dy1,
		path.currentPoint.x+dx2,
		path.currentPoint.y+dy2)
}

func (path *pathFixed) closePath() Status {
	var status Status
	if !path.hasCurrentPoint {
		return StatusSuccess
	}

	/*
	 * Add a line_to, to compute flags and solve any degeneracy.
	 * It will be removed later (if it was actually added).
	 */
	status = path.lineTo(path.lastMovePoint.x, path.lastMovePoint.y)
	if status != 0 {
		return status
	}

	/*
	 * If the command used to close the path is a line_to, drop it.
	 * We must check that last command is actually a line_to,
	 * because the path could have been closed with a curve_to (and
	 * the previous line_to not added as it would be degenerate).
	 */
	if path.lastOp() == pathOpLineTo {
		path.dropLineTo()
	}

	path.needsMoveTo = true /* After close_path, add an implicit move_to */

	return path.add(pathOpClosePath, nil)
}

func (path *pathFixed) getCurrentPoint() (x, y fixed, ok bool) {
	if !path.hasCurrentPoint {
		return 0, 0, false
	}

	x = path.currentPoint.x
	y = path.currentPoint.y
	return x, y, true
}

func (path *pathFixed) add(op pathOp, points []point) Status {
	buf := path.tail()
	if len(buf.ops)+1 > cap(buf.ops) ||
		len(buf.points)+len(points) > cap(buf.points) {
		buf = pathBufCreate(len(buf.ops)*2, len(buf.points)*2)
		path.addBuf(buf)
	}

	buf.addOp(op)
	buf.addPoints(points)
	return StatusSuccess
}

func (path *pathFixed) addBuf(buf *pathBuf) {
	path.buf = append(path.buf, buf)
}

func pathBufCreate(numOps, numPoints int) *pathBuf {
	return &pathBuf{
		ops:    make([]pathOp, 0, numOps),
		points: make([]point, 0, numPoints),
	}
}

func (buf *pathBuf) addOp(op pathOp) {
	buf.ops = append(buf.ops, op)
}

func (buf *pathBuf) addPoints(points []point) {
	buf.points = append(buf.points, points...)
}

func (path *pathFixed) interpret(moveTo pathFixedMoveToFunc,
	lineTo pathFixedLineToFunc,
	curveTo pathFixedCurveToFunc,
	closePath pathFixedClosePathFunc,
	closure interface{}) Status {

	var status Status

	var points []point

	for _, buf := range path.buf {
		points = buf.points
		for _, op := range buf.ops {
			switch op {
			case pathOpMoveTo:
				status = moveTo(closure, &points[0])
				points = points[1:]
			case pathOpLineTo:
				status = lineTo(closure, &points[0])
				points = points[1:]

			case pathOpCurveTo:
				status = curveTo(closure, &points[0], &points[1], &points[2])
				points = points[3:]

			case pathOpClosePath:
				status = closePath(closure)

			default:
				assertNotReached()
			}
			if status != 0 {
				return status
			}
		}
	}

	if path.needsMoveTo && path.hasCurrentPoint {
		return moveTo(closure, &path.currentPoint)
	}

	return StatusSuccess
}

type pathFixedAppendClosure struct {
	offset point
	path   *pathFixed
}

func appendMoveTo(abstractClosure interface{}, point *point) Status {
	closure := abstractClosure.(*pathFixedAppendClosure)
	return closure.path.moveTo(point.x+closure.offset.x,
		point.y+closure.offset.y)
}

func appendLineTo(abstractClosure interface{}, point *point) Status {
	closure := abstractClosure.(*pathFixedAppendClosure)
	return closure.path.lineTo(point.x+closure.offset.x,
		point.y+closure.offset.y)
}

func appendCurveTo(abstractClosure interface{}, p0, p1, p2 *point) Status {
	closure := abstractClosure.(*pathFixedAppendClosure)
	return closure.path.curveTo(p0.x+closure.offset.x,
		p0.y+closure.offset.y,
		p1.x+closure.offset.x,
		p1.y+closure.offset.y,
		p2.x+closure.offset.x,
		p2.y+closure.offset.y)
}

func appendClosePath(abstractClosure interface{}) Status {
	closure := abstractClosure.(*pathFixedAppendClosure)
	return closure.path.closePath()
}

func (path *pathFixed) append(other *pathFixed, tx, ty fixed) Status {
	var closure pathFixedAppendClosure
	closure.path = path
	closure.offset.x = tx
	closure.offset.y = ty

	return path.interpret(appendMoveTo, appendLineTo, appendCurveTo,
		appendClosePath, &closure)
}

func (path *pathFixed) offsetAndScale(offX, offY, scaleX, scaleY fixed) {
	if scaleX == fixedOne && scaleY == fixedOne {
		path.translate(offX, offY)
		return
	}
	path.lastMovePoint.x = scaleX.mul(path.lastMovePoint.x) + offX
	path.lastMovePoint.y = scaleY.mul(path.lastMovePoint.y) + offY
	path.currentPoint.x = scaleX.mul(path.currentPoint.x) + offX
	path.currentPoint.y = scaleY.mul(path.currentPoint.y) + offY

	path.fillMaybeRegion0 = true

	for i := 0; i < len(path.points); i++ {
		if scaleX != fixedOne {
			path.points[i].x = path.points[i].x.mul(scaleX)
		}
		path.points[i].x += offX

		if scaleY != fixedOne {
			path.points[i].y = path.points[i].y.mul(scaleY)
		}
		path.points[i].y += offY

		if path.fillMaybeRegion0 {
			path.fillMaybeRegion0 = path.points[i].x.IsInteger() &&
				path.points[i].y.IsInteger()
		}
	}
	path.fillMaybeRegion0 = path.fillMaybeRegion0 && path.fillIsRectilinear0
	path.extents.p1.x = scaleX.mul(path.extents.p1.x) + offX
	path.extents.p2.x = scaleX.mul(path.extents.p2.x) + offX
	if scaleX < 0 {
		t := path.extents.p1.x
		path.extents.p1.x = path.extents.p2.x
		path.extents.p2.x = t
	}

	path.extents.p1.y = scaleY.mul(path.extents.p1.y) + offY
	path.extents.p2.y = scaleY.mul(path.extents.p2.y) + offY
	if scaleY < 0 {
		t := path.extents.p1.y
		path.extents.p1.y = path.extents.p2.y
		path.extents.p2.y = t
	}
}

func (path *pathFixed) translate(offX, offY fixed) {
	if offX == 0 && offY == 0 {
		return
	}

	path.lastMovePoint.x += offX
	path.lastMovePoint.y += offY
	path.currentPoint.x += offX
	path.currentPoint.y += offY

	path.fillMaybeRegion0 = true

	for i := 0; i < len(path.points); i++ {
		path.points[i].x += offX
		path.points[i].y += offY

		if path.fillMaybeRegion0 {
			path.fillMaybeRegion0 = path.points[i].x.IsInteger() &&
				path.points[i].y.IsInteger()
		}
	}
	path.fillMaybeRegion0 = path.fillMaybeRegion0 && path.fillIsRectilinear0

	path.extents.p1.x += offX
	path.extents.p1.y += offY
	path.extents.p2.x += offX
	path.extents.p2.y += offY
}

func pathFixedTransformPoint(p *point, matrix *Matrix) {
	dx := p.x.toDouble()
	dy := p.y.toDouble()
	matrix.TransformPoint(&dx, &dy)
	p.x = fixedFromDouble(dx)
	p.y = fixedFromDouble(dy)
}

func (path *pathFixed) transform(matrix *Matrix) {
	if matrix.YX == 0 && matrix.XY == 0 {
		/* Fast path for the common case of scale+transform */
		path.offsetAndScale(fixedFromDouble(matrix.X0),
			fixedFromDouble(matrix.Y0),
			fixedFromDouble(matrix.XX),
			fixedFromDouble(matrix.YY))
		return
	}

	var extents box
	var point point

	pathFixedTransformPoint(&path.lastMovePoint, matrix)
	pathFixedTransformPoint(&path.currentPoint, matrix)

	if len(path.points) == 0 {
		return
	}

	extents = path.extents
	point = path.points[0]
	pathFixedTransformPoint(&point, matrix)
	path.extents.set(&point, &point)

	for i := 0; i < len(path.points); i++ {
		pathFixedTransformPoint(&path.points[i], matrix)
		path.extents.addPoint(&path.points[i])
	}

	if path.hasCurveTo {
		var isTight bool
		matrix.transformBoundingBoxFixed(&extents, &isTight)
		if !isTight {
			var hasExtents bool
			hasExtents = pathBounderExtents(path, &extents)
			if !hasExtents {
				panic("assert failed hasExtents == true")
			}
		}
		path.extents = extents
	}

	/* flags might become more strict than needed */
	path.strokeIsRectilinear0 = false
	path.fillIsRectilinear0 = false
	path.fillIsEmpty0 = false
	path.fillMaybeRegion0 = false
}

/* Closure for path flattening */
type pathFlattener struct {
	tolerance    float64
	currentPoint point
	moveTo       pathFixedMoveToFunc
	lineTo       pathFixedLineToFunc
	closePath    pathFixedClosePathFunc
	closure      interface{}
}

func cpfMoveTo(closure interface{}, point *point) Status {
	pf := closure.(*pathFlattener)
	pf.currentPoint = *point
	return pf.moveTo(pf.closure, point)
}

func cpfLineTo(closure interface{}, point *point) Status {
	pf := closure.(*pathFlattener)
	pf.currentPoint = *point
	return pf.lineTo(pf.closure, point)
}

func cpfCurveTo(closure interface{}, p1, p2, p3 *point) Status {
	pf := closure.(*pathFlattener)
	var spline spline
	p0 := &pf.currentPoint

	splineAddPoint := func(closure interface{}, point *point, tangent *slope) Status {
		return pf.lineTo(closure, point)
	}
	if !spline.init(splineAddPoint, pf.closure, p0, p1, p2, p3) {
		return cpfLineTo(closure, p3)
	}

	pf.currentPoint = *p3
	return spline.decompose(pf.tolerance)
}

func cpfClosePath(closure interface{}) Status {
	pf := closure.(*pathFlattener)
	return pf.closePath(pf.closure)
}

func (path *pathFixed) interpretFlat(moveTo pathFixedMoveToFunc, lineTo pathFixedLineToFunc,
	closePath pathFixedClosePathFunc, closure interface{}, tolerance float64) Status {
	if !path.hasCurveTo {
		return path.interpret(moveTo, lineTo, nil, closePath, closure)
	}

	var flattener pathFlattener
	flattener.tolerance = tolerance
	flattener.moveTo = moveTo
	flattener.lineTo = lineTo
	flattener.closePath = closePath
	flattener.closure = closure
	return path.interpret(cpfMoveTo, cpfLineTo, cpfCurveTo, cpfClosePath, &flattener)
}

func canonicalBox(box *box, p1, p2 *point) {
	if p1.x <= p2.x {
		box.p1.x = p1.x
		box.p2.x = p2.x
	} else {
		box.p1.x = p2.x
		box.p2.x = p1.x
	}

	if p1.y <= p2.y {
		box.p1.y = p1.y
		box.p2.y = p2.y
	} else {
		box.p1.y = p2.y
		box.p2.y = p1.y
	}
}

func (path *pathFixed) isQuad() bool {
	/* Do we have the right number of ops? */
	if len(path.ops) < 4 || len(path.ops) > 6 {
		return false
	}
	//len(path.ops) is 4,5,6

	/* Check whether the ops are those that would be used for a rectangle */
	if path.ops[0] != pathOpMoveTo ||
		path.ops[1] != pathOpLineTo ||
		path.ops[2] != pathOpLineTo ||
		path.ops[3] != pathOpLineTo {
		return false
	}

	/* we accept an implicit close for filled paths */
	if len(path.ops) > 4 {
		/* Now, there are choices. The rectangle might end with a LINE_TO
		 * (to the original point), but this isn't required. If it
		 * doesn't, then it must end with a CLOSE_PATH. */
		if path.ops[4] == pathOpLineTo {
			if path.points[4].x != path.points[0].x ||
				path.points[4].y != path.points[0].y {
				return false
			}
		} else if path.ops[4] != pathOpClosePath {
			return false
		}

		if len(path.ops) == 6 {
			/* A trailing CLOSE_PATH or MOVE_TO is ok */
			if path.ops[5] != pathOpMoveTo &&
				path.ops[5] != pathOpClosePath {
				return false
			}
		}
	}
	return true
}

func pointsFromRect(points []point) bool {
	if points[0].y == points[1].y &&
		points[1].x == points[2].x &&
		points[2].y == points[3].y &&
		points[3].x == points[0].x {
		return true
	}
	if points[0].x == points[1].x &&
		points[1].y == points[2].y &&
		points[2].x == points[3].x &&
		points[3].y == points[0].y {
		return true
	}
	return false
}

func (path *pathFixed) isBox(box *box) bool {
	if !path.fillIsRectilinear0 {
		return false
	}

	if !path.isQuad() {
		return false
	}

	if pointsFromRect(path.points) {
		canonicalBox(box, &path.points[0], &path.points[2])
		return true
	}

	return false
}

/* Determine whether two lines A->B and C->D intersect based on the
 * algorithm described here: http://paulbourke.net/geometry/pointlineplane/ */
func linesIntersectOrAreCoincident(a, b, c, d point) bool {
	denominator := (int64(d.y-c.y) * int64(b.x-a.x)) -
		(int64(d.x-c.x) * int64(b.y-a.y))
	numeratorA := (int64(d.x-c.x) * int64(a.y-c.y)) -
		(int64(d.y-c.y) * int64(a.x-c.x))
	numeratorB := (int64(b.x-a.x) * int64(a.y-c.y)) -
		(int64(b.y-a.y) * int64(a.x-c.x))

	if denominator == 0 {
		/* If the denominator and numerators are both zero,
		 * the lines are coincident. */
		if numeratorA == 0 && numeratorB == 0 {
			return true
		}

		/* Otherwise, a zero denominator indicates the lines are
		*  parallel and never intersect. */
		return false
	}

	/* The lines intersect if both quotients are between 0 and 1 (exclusive). */

	/* We first test whether either quotient is a negative number. */
	denominatorNegative := denominator < 0
	if (numeratorA < 0) != denominatorNegative {
		return false
	}
	if (numeratorB < 0) != denominatorNegative {
		return false
	}

	/* A zero quotient indicates an "intersection" at an endpoint, which
	 * we aren't considering a true intersection. */
	if numeratorA == 0 || numeratorB == 0 {
		return false
	}

	/* If the absolute value of the numerator is larger than or equal to the
	 * denominator the result of the division would be greater than or equal
	 * to one. */
	if !denominatorNegative {
		if !(numeratorA < denominator) || !(numeratorB < denominator) {
			return false
		}
	} else {
		if !(denominator < numeratorA) || !(denominator < numeratorB) {
			return false
		}
	}

	return true
}

func (path *pathFixed) isSimpleQuad() bool {
	if !path.isQuad() {
		return false
	}

	points := path.points
	if pointsFromRect(points) {
		return true
	}

	if linesIntersectOrAreCoincident(points[0], points[1], points[3], points[2]) {
		return false
	}

	if linesIntersectOrAreCoincident(points[0], points[3], points[1], points[2]) {
		return false
	}

	return true
}

func (path *pathFixed) isStrokeBox(box *box) bool {
	if !path.fillIsRectilinear0 {
		return false
	}

	/* Do we have the right number of ops? */
	if len(path.ops) != 5 {
		return false
	}

	/* Check whether the ops are those that would be used for a rectangle */
	if path.ops[0] != pathOpMoveTo ||
		path.ops[1] != pathOpLineTo ||
		path.ops[2] != pathOpLineTo ||
		path.ops[3] != pathOpLineTo ||
		path.ops[4] != pathOpClosePath {
		return false
	}

	/* Ok, we may have a box, if the points line up */
	if path.points[0].y == path.points[1].y &&
		path.points[1].x == path.points[2].x &&
		path.points[2].y == path.points[3].y &&
		path.points[3].x == path.points[0].x {

		canonicalBox(box, &path.points[0], &path.points[2])
		return true
	}

	if path.points[0].x == path.points[1].x &&
		path.points[1].y == path.points[2].y &&
		path.points[2].x == path.points[3].x &&
		path.points[3].y == path.points[0].y {

		canonicalBox(box, &path.points[0], &path.points[2])
		return true
	}

	return false
}

func (path *pathFixed) isRectangle(box *box) bool {
	if !path.isBox(box) {
		return false
	}

	/* This check is valid because the current implementation of
	 * _cairo_path_fixed_is_box () only accepts rectangles like:
	 * move,line,line,line[,line|close[,close|move]]. */
	if len(path.ops) > 4 {
		return true
	}
	return false
}
