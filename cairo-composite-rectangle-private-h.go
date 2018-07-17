package cairo

type compoisteRectangles struct {
	surface Surface
	op      Operator

	source, mask, destionation RectangleInt

	bounded   RectangleInt
	unbounded RectangleInt
	isBounded bool

	sourceSampleArea, maskSampleArea RectangleInt

	sourcePattern patternUnion
	maskPattern   patternUnion

	originalSourcePattern Pattern
	originalMaskPattern   Pattern
}
