package cairo

import "container/list"

type patternObserver struct {
	notify func(po *patternObserver, pattern *Pattern, flags uint)
	link   *list.List
}

const (
	patternNotifyMatrix  = 0x1
	patternNotifyFilter  = 0x2
	patternNotifyExtend  = 0x4
	patternNotifyOpacity = 0x9
)

type Pattern struct {
	status            Status
	observers         *list.List
	type0             PatternType
	filter            Filter
	extend            Extend
	hasComponentAlpha bool
	matrix            Matrix
	opacity           float64
}

type solidPattern struct {
	base  Pattern
	color color
}

type surfacePattern struct {
	base    Pattern
	surface *Surface
}

type gradientStop struct {
	offset float64
	color  colorStop
}

type gradientPattern struct {
	base         Pattern
	nStops       uint
	stopsSize    uint
	stops        []gradientStop
	stopEmbedded [2]gradientStop
}

type linearPattern struct {
	base     gradientPattern
	pd1, pd2 float64
}

type radialPattern struct {
	base     gradientPattern
	cd1, cd2 circleDouble
}

type gradientPatternUnion struct {
	base   gradientPattern
	linear linearPattern
	radial radialPattern
}

type meshPatch struct {
	points [4][4]pointDouble // or [16]pointDouble?
	colors [4]color
}

type meshPattern struct {
	base            Pattern
	patches         []meshPatch
	currentPatch    *meshPatch
	currentSize     int
	hasControlPoint [4]bool
	hasColor        [4]bool
}

type rasterSourcePattern struct {
	base     Pattern
	content  Content
	extents  RectangleInt
	acquire  func()
	release  func()
	snapshot func()
	copy     func()
	finish   func()
}

type patternUnion struct {
	base         Pattern
	solid        solidPattern
	surface      surfacePattern
	gradient     gradientPatternUnion
	mesh         meshPattern
	rasterSource rasterSourcePattern
}
