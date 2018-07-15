package cairo

type edgeT struct {
	next, prev, right *edgeT
	x, top            fixed
	dir               int
}

type rectangleT struct {
	left, right edgeT
	top, bottom int32
}

func pqParentIndex(i int) int {
	return i >> 1
}

const pqFirstEntry = 1

func pqLeftChildIndex(i int) int {
	return i << 1
}

type sweepLine struct {
	rectangles     []*rectangleT
	stop           []*rectangleT
	head, tail     edgeT
	insert, cursor *edgeT
	currentY       int32
	lastY          int32
	insertX        int32
	fillRule       FillRule
	doTraps        bool
	container      interface{}
	unwind         jmpBuf
}

type jmpBuf struct {
}
