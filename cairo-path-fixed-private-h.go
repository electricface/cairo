package cairo

import (
	"container/list"
	"unsafe"
)

type pathOpType byte

const (
	pathOpMoveTo pathOpType = iota
	pathOpLineTo
	pathOpCurveTo
	pathOpClosePath
)

type pathOp byte

type pathBuf struct {
	link   *list.List
	ops    []pathOp
	points []point
}

type pathBufFixed struct {
	base   pathBuf
	ops    [pathBufSize]pathOp
	points [pathBufSize]point
}

const pathBufSize = 512 - unsafe.Sizeof(pathBuf{})/
	(2*unsafe.Sizeof(point{})+unsafe.Sizeof(pathOp(0)))
