package cairo

type gstate struct {
	op                 Operator
	opacity            float64
	tolerance          float64
	antialias          Antialias
	strokeStyle        strokeStyle
	fillRule           FillRule
	fontFace           *fontFace
	scaledFont         *scaledFont
	previousScaledFont *scaledFont
	fontMatrix         Matrix
	fontOptions        fontOptions
	clip               *clip
	target             *Surface
	parentTarget       *Surface
	originalTarget     *Surface

	deviceTransformObserver observer

	ctm              Matrix
	ctmInverse       Matrix
	sourceCtmInverse Matrix
	isIdentity       bool
	source           *Pattern
	next             *gstate
}
