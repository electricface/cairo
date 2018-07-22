package cairo

import (
	"unsafe"
)

const (
	pathOpMoveTo pathOp = iota
	pathOpLineTo
	pathOpCurveTo
	pathOpClosePath
)

type pathOp byte

type pathBuf struct {
	ops    []pathOp
	points []point
}

type pathBufFixed struct {
	ops    [pathBufSize]pathOp
	points [2 * pathBufSize]point
}

const pathBufSize = 512 - unsafe.Sizeof(pathBuf{})/
	(2*unsafe.Sizeof(point{})+unsafe.Sizeof(pathOp(0)))

type pathFixed struct {
	lastMovePoint        point
	currentPoint         point
	hasCurrentPoint      bool
	needsMoveTo          bool
	hasExtents           bool
	hasCurveTo           bool
	strokeIsRectilinear0 bool
	fillIsRectilinear0   bool
	fillMaybeRegion0     bool
	fillIsEmpty0         bool
	extents              box
	buf0                 pathBufFixed
	buf                  []*pathBuf
}

//type pathFixedIter struct {
//	first       *pathBuf
//	buf         *pathBuf
//	nOp, nPoint uint
//}

func (path *pathFixed) fillIsEmpty() bool {
	return path.fillIsEmpty0
}

func (path *pathFixed) fillIsRectilinear() bool {
	if !path.fillIsRectilinear0 {
		return false
	}

	if !path.hasCurrentPoint || path.needsMoveTo {
		return true
	}

	/* check whether the implicit close preserves the rectilinear property */
	return path.currentPoint.x == path.lastMovePoint.x ||
		path.currentPoint.y == path.lastMovePoint.y
}

func (path *pathFixed) strokeIsRectilinear() bool {
	return path.strokeIsRectilinear0
}

func (path *pathFixed) fillMaybeRegion() bool {
	if !path.fillMaybeRegion0 {
		return false
	}

	if !path.hasCurrentPoint || path.needsMoveTo {
		return true
	}

	/* check whether the implicit close preserves the rectilinear property
	 * (the integer point property is automatically preserved)
	 */
	return path.currentPoint.x == path.lastMovePoint.x ||
		path.currentPoint.y == path.lastMovePoint.y
}

//func (path *pathFixed) head() *pathBuf {
//	return &path.buf.base
//}
//
//func (pb *pathBuf) next() *pathBuf {
//	pb.link.next
//}
//
//func (pb *pathBuf) prev() *pathBuf {
//
//}
