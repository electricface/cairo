package cairo

import "container/list"

type boxes struct {
	status         Status
	limit          box
	limits         []box
	numBoxes       int
	isPixelAligned bool

	chunks *list.List // elem type boxesChunk
}

type boxesChunk struct {
	boxes []box
}

func makeBoxesChunk(cap int) boxesChunk {
	return boxesChunk{
		make([]box, 0, cap),
	}
}
