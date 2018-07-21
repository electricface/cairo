package cairo

import (
	"container/list"
	"math"
)

type array struct {
}

type cache struct {
}

type compositeRectangles struct {
}

type color struct {
	red, green, blue, alpha                     float64
	redShort, greenShort, blueShort, alphaShort uint16
}

type colorStop struct {
	/* unpremultiplied */
	red, green, blue, alpha float64
	/* unpremultipled, for convenience */
	redShort, greenShort, blueShort, alphaShort uint16
}

type contour struct {
}

type contourChain struct {
}

type contourIter struct {
}

type damage struct {
}

type deviceBackend struct {
}

type gstate struct {
}

type gstateBackend struct {
}

type imageSurface struct {
}

type observer struct {
	link     list.List
	callback func(self *observer, arg interface{})
}

type outputStream struct {
}

type paginatedSurfaceBackend struct {
}

type glyphSize struct {
}

type scaledFontSubsets struct {
}

type surfaceBackend struct {
}

type surfaceObserver struct {
}

type surfaceSnapshot struct {
}

type surfaceSubsurface struct {
}

type surfaceWrapper struct {
}

type tristrip struct {
}

type xlibScreenInfo struct {
}

type userDataArray struct {
}

type scaledFontPrivate struct {
}

type scaledGlyph struct {
}

type scaledGlyphPrivate struct {
}

type compositor struct {
}

type fallbackCompositor struct {
}

type maskCompositor struct {
}

type trapsCompositor struct {
}

type spansCompositor struct {
}

type lcdFilter int

const (
	lcdFilterDefault lcdFilter = iota
	lcdFilterNone
	lcdFilterIntraPixel
	lcdFilterFIR3
	lcdFilterFIR5
)

type roundGlyphPositions int

const (
	roundGlyphPosDefault roundGlyphPositions = iota
	roundGlyphPosOn
	roundGlyphPosOff
)

type fontOptions struct {
	antialias           Antialias
	subpixelOrder       SubpixelOrder
	lcdFilter           lcdFilter
	hintStyle           HintStyle
	hintMetrics         HintMetrics
	roundGlyphPositions roundGlyphPositions
	variations          string
}

type glyphTextInfo struct {
	utf8         string
	textClusters []TextCluster
	clusterFlag  TextClusterFlag
}

type paginatedMode int

const (
	paginatedModeAnalyze paginatedMode = iota
	paginatedModeRender
	paginatedModeFallback
)

type internalSurfaceType int

const (
	internalSurfaceTypeSnapshot internalDeviceType = iota + 0x1000
	internalSurfaceTypePaginated
	internalSurfaceTypeAnalysis
	internalSurfaceTypeObserver
	internalSurfaceTypeTestFallback
	internalSurfaceTypeTestPaginated

	internalSurfaceTypeWrapping
	internalSurfaceTypeNull
	internalSurfaceTypeType3Glyph
)

type internalDeviceType int

const (
	internalDeviceTypeObserver internalDeviceType = 0x1000
)

const hasTestPaginatedSurface = true

type slope struct {
	dx, dy fixed
}

type distance slope

type pointDouble struct {
	x, y float64
}

type circleDouble struct {
	center pointDouble
	radius float64
}

type distanceDouble struct {
	dx, dy float64
}

type boxDouble struct {
	p1, p2 pointDouble
}

type line struct {
	p1, p2 point
}

type box line

type trapzoid struct {
	top, bottom fixed
	left, right line
}

type pointInt struct {
	x, y int
}

const intMax = math.MaxInt32
const intMin = math.MinInt32

const rectIntMin = intMin >> fixedFracBits
const rectIntMax = intMax >> fixedFracBits

type direction int

const (
	directionForward direction = iota
	directionReverse
)

type edge struct {
	line        line
	top, bottom int
	dir         int
}

type polygon struct {
	status  Status
	extents box
	limit   box
	limits  []box

	edges         []edge
	edgesSize     int
	edgesEmbedded [32]edge
}

type splineKnots struct {
	a, b, c, d point
}

type splineAddPointFunc func(closure interface{}, point *point, tangent *slope) Status

type spline struct {
	addPointFunc splineAddPointFunc
	closure      interface{}

	knots        splineKnots
	initialSlope slope
	finalSlope   slope
	hasPoint     bool
	lastPoint    point
}

type penVertex struct {
	point   point
	slopCcw slope
	slopCw  slope
}

type pend struct {
	radius           float64
	tolerance        float64
	vertices         []penVertex
	verticesEmbedded [32]penVertex
}

type strokeStyle struct {
	lineWidth  float64
	lineCap    LineCap
	lineJoin   LineJoin
	miterLimit float64
	dashes     []float64
	dashOffset float64
}

type formatMasks struct {
	bpp                                     int
	alphaMask, redMask, greenMask, blueMask uint
}

type stock int

const (
	stockWhite stock = iota
	stockBlack
	stockTransparent
	stockNumColors
)

type imageTransparency int

const (
	imageIsOpaque imageTransparency = iota
	imageHasBilevelAlpha
	imageHasAlpha
	imageUnknown
)

type imageColor int

const (
	imageIsColor imageColor = iota
	imageIsGrayscale
	imageIsMonoChrome
	imageUnknownColor
)

type mimeData struct {
	data []byte
}

type unscaledFont struct {
	// hashEntry
	backend *unscaledFontBackend
}
