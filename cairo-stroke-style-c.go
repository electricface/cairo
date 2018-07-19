package cairo

func (style *strokeStyle) init() {
	style.lineWidth = gstateLineWidthDefault
	style.lineCap = gstateLineCapDefault
	style.lineJoin = gstateLineJoinDefault
	style.miterLimit = gstateMiterLimitDefault
	style.dash = nil
	style.dashOffset = 0
}

func (style *strokeStyle) initCopy(other *strokeStyle) Status {
	style.lineWidth = other.lineWidth
	style.lineCap = other.lineCap
	style.lineJoin = other.lineJoin
	style.miterLimit = other.miterLimit

	if other.dash == nil {
		style.dash = nil
	} else {
		style.dash = make([]float64, len(other.dash))
		copy(style.dash, other.dash)
	}

	style.dashOffset = other.dashOffset
	return StatusSuccess
}

func (style *strokeStyle) fini() {
	style.dash = nil
}

func (style *strokeStyle) maxDistanceFromPath(path *pathFixed, ctm *Matrix, dx, dy *float64) {
	styleExpansion := 0.5
	if style.lineCap == LineCapSquare {
		styleExpansion = mSqrt1_2
	}

	if style.lineJoin == LineJoinMiter &&
}
