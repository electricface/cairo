package cairo

import "math"

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

func (style *strokeStyle) maxDistanceFromPath(path *pathFixed, ctm *Matrix) (dx, dy float64) {
	styleExpansion := 0.5
	if style.lineCap == LineCapSquare {
		styleExpansion = mSqrt1_2
	}

	if style.lineJoin == LineJoinMiter &&
		!path.strokeIsRectilinear &&
		styleExpansion < math.Sqrt2*style.miterLimit {
		styleExpansion = math.Sqrt2 * style.miterLimit
	}
	styleExpansion *= style.lineWidth

	if ctm.hasUnityScale() {
		dx = styleExpansion
		dy = styleExpansion
	} else {
		dx = styleExpansion * math.Hypot(ctm.XX, ctm.XY)
		dy = styleExpansion * math.Hypot(ctm.YY, ctm.YX)
	}
	return
}

func (style *strokeStyle) maxLineDistanceFromPath(path *pathFixed, ctm *Matrix) (dx, dy float64) {
	styleExpansion := 0.5 * style.lineWidth
	if ctm.hasUnityScale() {
		dx = styleExpansion
		dy = styleExpansion
	} else {
		dx = styleExpansion * math.Hypot(ctm.XX, ctm.XY)
		dy = styleExpansion * math.Hypot(ctm.YY, ctm.YX)
	}
	return
}
