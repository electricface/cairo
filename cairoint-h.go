package cairo

import (
	"math"
	"unsafe"
)

const mSqrt1_2 = 0.707106781186547524400844362104849039

func notReached() {
	panic("not reached")
}

func alphaIsClear(alpha float64) bool {
	return alpha <= (float64(0x00ff) / float64(0xffff))
}

func alphaShortIsClear(alpha uint16) bool {
	return alpha <= 0x00ff
}

func alphaIsOpaque(alpha float64) bool {
	return alpha >= (float64(0xff00) / float64(0xffff))
}

func alphaShortIsOpaque(alpha uint16) bool {
	return alpha >= 0xff00
}

func alphaIsZero(alpha float64) bool {
	return alpha <= 0.0
}

func (c color) isClear() bool {
	return alphaShortIsClear(c.alphaShort)
}

func (c color) isOpaque() bool {
	return alphaShortIsOpaque(c.alphaShort)
}

func popCount(mask uint32) int {
	var y int
	y = int(mask>>1) & 033333333333
	y = int(mask) - y - ((y >> 1) & 033333333333)
	return ((y + (y >> 3)) & 030707070707) % 077
}

func isLittleEndian() bool {
	var i = 1
	return *(*byte)(unsafe.Pointer(&i)) == 1
}

type fontFace struct {
	status  Status
	backend *fontFaceBackend
}

type unscaledFontBackend interface {
	destory(font *unscaledFont)
}

type toyFontFace struct {
	base       fontFace
	family     string
	ownsFamily bool
	slant      FontSlant
	weight     FontWeight
	implFace   *fontFace
}

type scaledGlyphInfo int

const (
	scaledGlyphInfoMetrics = 1 << iota
	scaledGlyphInfoSurface
	scaledGlyphInfoPath
	scaledGlyphInfoRecordingSurface
	scaledGlyphInfoColorSurface
)

type scaledFontSubset struct {
	scaledFont *scaledFont
	fontId     uint
	subsetId   uint

	glyphs                  []uint64
	utf8                    []string
	glyphNames              []string
	toLatinChar             []int
	latinToSubsetGlyphIndex []uint64
	numGlyphs               uint
	isComposite             bool
	isScaled                bool
	isLatin                 bool
}

type scaledFontBackend interface {
	fini()
	scaledGlyphInit(scaledGlyph scaledGlyph, info scaledGlyphInfo)
	textToGlphs()
	ucs4ToIndex()
	loadTruetypeTable()
	indexToUcs4()
	isSynthetic() (bool, Status)
	indexToGlyphName()
	loadType1Data()
	hasColorGlyphs()
}

type fontFaceBackend interface {
	getType()
	setType()
	createForToy()
	destroy()
	scaledFontCreate()
	getImplementation()
}

type surfaceAttributes struct {
	matrix            Matrix
	extend            Extend
	filter            Filter
	hasComponentAlpha bool
	xOffset, yOffset  bool
	extra             interface{}
}

type strokeFace struct {
	ccw       point
	point     point
	cw        point
	devVector slope
	devSlope  pointDouble
	usrVector pointDouble
	length    float64
}

func restrictValue(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func round(r float64) float64 {
	return math.Floor(r + 0.5)
}

func lRound(r float64) float64 {
	return round(r)
}

const (
	operatorBoundByMask    = 1 << 1
	opearatorBoundBySource = 1 << 2
)

var colorWhite = stockColor(stockWhite)
var colorBlack = stockColor(stockBlack)
var colorTransparent = stockColor(stockTransparent)

func (p *polygon) isEmpty() bool {
	return len(p.edges) == 0 ||
		p.extents.p2.x <= p.extents.p1.x
}

static inline cairo_bool_t
_cairo_matrix_is_identity (const cairo_matrix_t *matrix)
{
return (matrix->xx == 1.0 && matrix->yx == 0.0 &&
matrix->xy == 0.0 && matrix->yy == 1.0 &&
matrix->x0 == 0.0 && matrix->y0 == 0.0);
}
