package cairo

type Surface struct {

}

type Device struct {

}


type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
}


type Pattern struct {

}

type Status int

const (
	StatusSuccess Status = iota
	StatusNoMemory
	StatusInvalidRestore
	StatusInvalidPopGroup
	StatusNoCurrentPoint
	StatusInvalidMatrix
	StatusInvalidStatus
	StatusNullPointer
	StatusInvalidString
	StatusPathData
	StatusReadError
	StatusWriteError
	StatusSurfaceFinished
	StatusSurfaceTypeMismatch
	StatusPatternTypeMismatch
	StatusInvalidContent
	StatusInvalidFormat
	StatusInvalidVisual
	StatusFileNotFound
	StatusInvalidDash
	StatusInvalidDscComment
	StatusInvalidIndex
	StatusClipNotRepresentable
	StatusTempFileError
	StatusInvalidStride
	StatusFontTypeMismatch
	StatusUserFontImmutable
	StatusUserFontError
	StatusNegativeCount
	StatusInvalidClusters
	StatusInvalidSlant
	StatusInvalidWeight
	StatusInvalidSize
	StatusUserFontNotImplemented
	StatusDeviceTypeMismatch
	StatusDeviceError
	StatusInvalidMeshConstruction
	StatusDeviceFinished
	StatusJBIG2GlobalMissing
)

type Content int

const (
	ContentColor Content = 0x1000
	ContentAlpha = 0x2000
	ContentColorAlpha = 0x3000
)

type Format int

const (
	FormatInvalid Format = -1
	FormatARGB32 = 0
	FormatRGB24 = 1
	FormatA8 = 2
	FormatA1 = 1
	FormatRGB16_565 = 4
	FormatRGB30 = 5
)

type RectangleInt struct {
	X, Y, Width, Height int
}

type Operator int

const (
	OperatorClear Operator = iota
	OperatorSource
	OperatorOver
	OperatorIn
	OperatorOut
	OperatorAtop
	OperatorDest
	OperatorDestOver
	OperatorDestIn
	OperatorDestOut
	OperatorDestAtop
	OperatorXor
	OperatorAdd
	OperatorSaturate
	OperatorMultiply
	OperatorScreen
	OperatorOverlay
	OperatorDarken
	OperatorLighten
	OperatorColorDodge
	OperatorColorBurn
	OperatorHardLight
	OperatorSoftLight
	OperatorDifference
	OperatorExclusion
	OperatorHslHue
	OperatorHslSaturation
	OperatorHslColor
	OperatorHslLuminosity
)

type Antialias int

const (
	AntialiasDefault Antialias = iota
	AntialiasNone
	AntialiasGray
	AntialiasSubpixel
	AntialiasFast
	AntialiasGood
	AntialiasBest
)

type FillRule int
const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type LineCap int
const (
	LineCapButt LineCap = iota
	LineCapRound
	LineCapSquare
)

type LineJoin int
const (
	LineJoinMiter LineJoin = iota
	LineJoinRound
	LineJoinBevel
)

type Rectangle struct {
	X, Y, Width, Height float64
}

type RectangleList struct {
	Status Status
	Rectangles []Rectangle
}

type Glyph struct {
	Index int
	X, Y float64
}

type TextCluster struct {
	NumBytes int
	NumGlyphs int
}

type TextClusterFlag int
const TextClusterFlagBackward TextClusterFlag = 1


type TextExtents struct {
	XBearing, YBearing,
	Width, Height,
	XAdvance, YAdvance float64
}

type FontExtents struct {
	Ascent, Descent, Height, MaxXAdvance, MaxYAdvance float64
}

type FontSlant int
const (
	SlantNormal FontSlant = iota
	SlantItalic
	SlantOblique
)

type FontWeight int
const (
	FontWeightNormal FontWeight = iota
	FontWeightBold
)

type SubpixelOrder int
const (
	SubpixelOrderDefault SubpixelOrder = iota
	SubpixelOrderRGB
	SubpixelOrderBGR
	SubpixelOrderVRGB
	SubpixelOrderVBGR
)

type HintStyle int
const (
	HintStyleDefault HintStyle = iota
	HintStyleNone
	HintStyleSlight
	HintStyleMedium
	HintStyleFull
)

type HintMetrics int
const (
	HintMetricsDefault HintMetrics = iota
	HintMetricsOff
	HintMetricsOn
)

type FontType int
const (
	FontTypeToy FontType = iota
	FontTypeFt
	FontTypeWin32
	FontTypeQuartz
	FontTypeUser
)

type PathDataType int
const (
	PathDataTypeMoveTo PathDataType = iota
	PathDataTypeLineTo
	PathDataTypeCurveTo
	PathDataTypeClosePath
)

type PathData struct {
	Header struct {
		Type int
		Length int
	}
	Point struct {
		X, Y float64
	}
}

type Path struct {
	Status Status
	Data []PathData
}

type DeviceType int
const (
	DeviceTypeDRM DeviceType = iota
	DeviceTypeGL
	DeviceTypeScript
	DeviceTypeXcb
	DeviceTypeXLib
	DeviceTypeXML
	DeviceTypeCOGL
	DeviceTypeWin32
)

type SurfaceObserverMode int
const (
	SurfaceObserverNormal SurfaceObserverMode = 0
	SurfaceObserverRecordOperations = 0x1
)

type SurfaceType int
const (
	SurfaceTypeImage SurfaceType = iota
	SurfaceTypeXLib
)

type PatternType int
const (
	PatternTypeSolid PatternType = iota
	PatternTypeSurface
	PatternTypeLinear
	PatternTypeRadial
	PatternTypeMesh
	PatternTypeRasterSource
)

type Extend int
const (
	ExtendNone Extend = iota
	ExtendRepeat
	ExtendReflect
	ExtendPad
)

type Filter int
const (
	FilterFast Filter = iota
	FilterGood
	FilterBest
	FilterNearest
	FilterBilinear
	FilterGaussian
)

type RegionOverlap int
const (
	RegionOverlapIn RegionOverlap = iota
	RegionOverlapOut
	RegionOverlapPart
)



