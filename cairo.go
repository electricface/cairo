package cairo

type Cairo struct {

}

type Surface struct {

}

type Device struct {

}


type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
}


type Pattern struct {

}

const (
	StatusSuccess = iota
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

const (
	ContentColor = 0x1000
	ContentAlpha = 0x2000
	ContentColorAlpha = 0x3000
)

const (
	FormatInvalid = -1
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

const (
	OperatorClear = iota
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

const (
	AntialiasDefault = iota
	AntialiasNone
	AntialiasGray
	AntialiasSubpixel
	AntialiasFast
	AntialiasGood
	AntialiasBest
)

const (
	FillRuleWinding = iota
	FillRuleEvenOdd
)

const (
	LineCapButt = iota
	LineCapRound
	LineCapSquare
)

const (
	LineJoinMiter = iota
	LineJoinRound
	LineJoinBevel
)

type Rectangle struct {
	X, Y, Width, Height float64
}

type RectangleList struct {
	Status int
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

const TextClusterFlagBackward = 1


type TextExtents struct {
	XBearing, YBearing,
	Width, Height,
	XAdvance, YAdvance float64
}

type FontExtents struct {
	Ascent, Descent, Height, MaxXAdvance, MaxYAdvance float64
}

const (
	SlantNormal = iota
	SlantItalic
	SlantOblique
)

const (
	FontWeightNormal = iota
	FontWeightBold
)

const (
	SubpixelOrderDefault = iota
	SubpixelOrderRGB
	SubpixelOrderBGR
	SubpixelOrderVRGB
	SubpixelOrderVBGR
)

const (
	HintStyleDefault = iota
	HintStyleNone
	HintStyleSlight
	HintStyleMedium
	HintStyleFull
)

const (
	HintMetricsDefault = iota
	HintMetricsOff
	HintMetricsOn
)

