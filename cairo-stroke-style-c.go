package cairo

import "math"

func (style *strokeStyle) init() {
	style.lineWidth = gstateLineWidthDefault
	style.lineCap = gstateLineCapDefault
	style.lineJoin = gstateLineJoinDefault
	style.miterLimit = gstateMiterLimitDefault
	style.dashes = nil
	style.dashOffset = 0
}

func (style *strokeStyle) initCopy(other *strokeStyle) Status {
	style.lineWidth = other.lineWidth
	style.lineCap = other.lineCap
	style.lineJoin = other.lineJoin
	style.miterLimit = other.miterLimit

	if other.dashes == nil {
		style.dashes = nil
	} else {
		style.dashes = make([]float64, len(other.dashes))
		copy(style.dashes, other.dashes)
	}

	style.dashOffset = other.dashOffset
	return StatusSuccess
}

func (style *strokeStyle) fini() {
	style.dashes = nil
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

func (style *strokeStyle) maxJoinDistanceFromPath(path *pathFixed, ctm *Matrix) (dx, dy float64) {
	styleExpansion := 0.5

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

func (style *strokeStyle) dashPeriod() float64 {
	var period float64
	for _, value := range style.dashes {
		period += value
	}

	if len(style.dashes)&1 != 0 {
		period *= 2.0
	}
	return period
}

const roundMinsqApproximation = 8.0 * math.Pi / 32.0

func (style *strokeStyle) dashStroked() float64 {
	var stroked, capScale float64

	switch style.lineCap {
	case LineCapButt:
		capScale = 0.0
	case LineCapRound:
		capScale = roundMinsqApproximation
	case LineCapSquare:
		capScale = 1.0
	default:
		assertNotReached()
	}

	if len(style.dashes)&1 != 0 {
		/* Each dash element is used both as on and as off. The order in which they are summed is
		* irrelevant, so sum the coverage of one dash element, taken both on and off at each iteration */
		for _, value := range style.dashes {
			stroked += value + capScale*math.Min(value, style.lineWidth)
		}
	} else {
		/* Even (0, 2, ...) dashes are on and simply counted for the coverage, odd dashes are off, thus
		 * their coverage is approximated based on the area covered by the caps of adjacent on dases. */
		for i := 0; i+1 < len(style.dashes); i += 2 {
			stroked += style.dashes[i] + capScale*math.Min(style.dashes[i+1], style.lineWidth)
		}
	}

	return stroked
}

func (style *strokeStyle) dashCanApproximate(ctm *Matrix, tolerance float64) bool {
	if len(style.dashes) == 0 {
		return false
	}
	period := style.dashPeriod()
	return ctm.transformedCircleMajorAxis(period) < tolerance
}

func (style *strokeStyle) dashApproximate(ctm *Matrix, tolerance float64) (dashOffset float64,
	dashes []float64) {

	var coverage, scale, offset float64
	on := true
	coverage = style.dashStroked() / style.dashPeriod()
	coverage = math.Min(coverage, 1.0)
	scale = tolerance / ctm.transformedCircleMajorAxis(1.0)

	/* We stop searching for a starting point as soon as the
	* offset reaches zero.  Otherwise when an initial dash
	* segment shrinks to zero it will be skipped over. */
	offset = style.dashOffset
	var i int
	for offset > 0.0 && offset >= style.dashes[i] {
		offset -= style.dashes[i]
		on = !on
		i++
		if i == len(style.dashes) {
			i = 0
		}
	}

	dashes = make([]float64, 2)
	switch style.lineCap {
	default:
		assertNotReached()
	case LineCapButt:
		/* Simplified formula (substituting 0 for cap_scale): */
		dashes[0] = scale * coverage

	case LineCapRound:
		dashes[0] = math.Max(scale*(coverage-roundMinsqApproximation)/
			(1.0-roundMinsqApproximation),
			scale*coverage-roundMinsqApproximation*style.lineWidth)

	case LineCapSquare:
		/*
		* Special attention is needed to handle the case cap_scale == 1 (since the first solution
		* is either indeterminate or -inf in this case). Since dash lengths are always >=0, using
		* 0 as first solution always leads to the correct solution.
		 */
		dashes[0] = math.Max(0.0, scale*coverage-style.lineWidth)
	}

	dashes[1] = scale - dashes[0]
	if on {
		dashOffset = 0.0
	} else {
		dashOffset = dashes[0]
	}

	return
}
