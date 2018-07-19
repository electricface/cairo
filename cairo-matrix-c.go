package cairo

import "math"

func (m *Matrix) InitIdentity() {
	m.Init(1, 0,
		0, 1,
		0, 0)
}

func (m *Matrix) Init(xx, yx, xy, yy, x0, y0 float64) {
	m.XX = xx
	m.YX = yx
	m.XY = xy
	m.YY = yy
	m.X0 = x0
	m.Y0 = y0
}

func (m *Matrix) getAffine(xx, yx, xy, yy, x0, y0 *float64) {
	*xx = m.XX
	*yx = m.YX
	*xy = m.XY
	*yy = m.YY

	if x0 != nil {
		*x0 = m.X0
	}
	if y0 != nil {
		*y0 = m.Y0
	}
}

func (m *Matrix) InitTranslate(tx, ty float64) {
	m.Init(1, 0,
		0, 1,
		tx, ty)
}

func (m *Matrix) translate(tx, ty float64) {
	var tmp Matrix
	tmp.InitTranslate(tx, ty)
	m.Multiply(&tmp, m)
}

func (m *Matrix) InitScale(sx, sy float64) {
	m.Init(sx, 0,
		0, sy,
		0, 0)
}

func (m *Matrix) Scale(sx, sy float64) {
	var tmp Matrix
	tmp.InitScale(sx, sy)
	m.Multiply(&tmp, m)
}

func (m *Matrix) InitRotate(radians float64) {
	s := math.Sin(radians)
	c := math.Cos(radians)

	m.Init(c, s,
		-s, c,
		0, 0)
}

func (m *Matrix) Rotate(radians float64) {
	var tmp Matrix
	tmp.InitRotate(radians)
	m.Multiply(&tmp, m)
}

func (a *Matrix) Multiply(b, result *Matrix) {
	var r Matrix
	r.XX = a.XX*b.XX + a.YX*b.XY
	r.YX = a.XX*b.YX + a.YX*b.YY

	r.XY = a.XY*b.XX + a.YY*b.XY
	r.YY = a.XY*b.YX + a.YY*b.YY

	r.X0 = a.X0*b.XX + a.Y0*b.XY + b.X0
	r.Y0 = a.X0*b.YX + a.Y0*b.YY + b.Y0

	*result = r
}

func (a *Matrix) multiply(b, r *Matrix) {
	r.XX = a.XX*b.XX + a.YX*b.XY
	r.YX = a.XX*b.YX + a.YX*b.YY

	r.XY = a.XY*b.XX + a.YY*b.XY
	r.YY = a.XY*b.YX + a.YY*b.YY

	r.X0 = a.X0*b.XX + a.Y0*b.XY + b.X0
	r.Y0 = a.X0*b.YX + a.Y0*b.YY + b.Y0
}

func (m *Matrix) TransformDistance(dx, dy *float64) {
	newX := m.XX*(*dx) + m.XY*(*dy)
	newY := m.YX*(*dx) + m.YY*(*dy)
	*dx = newX
	*dy = newY
}

func (m *Matrix) TransformPoint(x, y *float64) {
	m.TransformDistance(x, y)
	*x += m.X0
	*y += m.Y0
}

func (m *Matrix) transformBoundingBox(x1, y1, x2, y2 *float64, isTight *bool) {
	var quadX, quadY [4]float64

	if m.XY == 0. && m.YX == 0. {
		/* non-rotation/skew matrix, just map the two extreme points */

		if m.XX != 1. {
			quadX[0] = *x1 * m.XX
			quadX[1] = *x2 * m.XX
			if quadX[0] < quadX[1] {
				*x1 = quadX[0]
				*x2 = quadX[1]
			} else {
				*x1 = quadX[1]
				*x2 = quadX[0]
			}
		}
		if m.X0 != 0. {
			*x1 += m.X0
			*x2 += m.X0
		}

		if m.YY != 1. {
			quadY[0] = *y1 * m.YY
			quadY[1] = *y2 * m.YY
			if quadY[0] < quadY[1] {
				*y1 = quadY[0]
				*y2 = quadY[1]
			} else {
				*y1 = quadY[1]
				*y2 = quadY[0]
			}
		}
		if m.Y0 != 0. {
			*y1 += m.Y0
			*y2 += m.Y0
		}

		if isTight != nil {
			*isTight = true
		}

		return
	}

	/* general matrix */
	quadX[0] = *x1
	quadY[0] = *y1
	m.TransformPoint(&quadX[0], &quadY[0])

	quadX[1] = *x2
	quadY[1] = *y1
	m.TransformPoint(&quadX[1], &quadY[1])

	quadX[2] = *x1
	quadY[2] = *y2
	m.TransformPoint(&quadX[2], &quadY[2])

	quadX[3] = *x2
	quadY[3] = *y2
	m.TransformPoint(&quadX[3], &quadY[3])

	minX := quadX[0]
	maxX := quadX[0]

	minY := quadY[0]
	maxY := quadY[0]

	for i := 1; i < 4; i++ {
		if quadX[i] < minX {
			minX = quadX[i]
		}
		if quadX[i] > maxX {
			maxX = quadX[i]
		}

		if quadY[i] < minY {
			minY = quadY[i]
		}
		if quadY[i] > maxY {
			maxY = quadY[i]
		}
	}

	*x1 = minX
	*y1 = minY
	*x2 = maxX
	*y2 = maxY

	if isTight != nil {
		/* it's tight if and only if the four corner points form an axis-aligned
		   rectangle.
		   And that's true if and only if we can derive corners 0 and 3 from
		   corners 1 and 2 in one of two straightforward ways...
		   We could use a tolerance here but for now we'll fall back to FALSE in the case
		   of floating point error.
		*/
		*isTight =
			(quadX[1] == quadX[0] && quadY[1] == quadY[3] &&
				quadX[2] == quadX[3] && quadY[2] == quadY[0]) ||
				(quadX[1] == quadX[3] && quadY[1] == quadY[0] &&
					quadX[2] == quadX[0] && quadY[2] == quadY[3])
	}
}

func (m *Matrix) transformBoundingBoxFixed(bBox *box, isTight *bool) {
	x1, y1, x2, y2 := bBox.toDoubles()
	m.transformBoundingBox(&x1, &y1, &x2, &y2, isTight)
	bBox.fromDoubles(x1, y1, x2, y2)
}

func (m *Matrix) scalarMultiply(scalar float64) {
	m.XX *= scalar
	m.YX *= scalar

	m.XY *= scalar
	m.YY *= scalar

	m.X0 *= scalar
	m.Y0 *= scalar
}

func (m *Matrix) computeAdjoint() {
	var a, b, c, d, tx, ty float64
	m.getAffine(&a, &b,
		&c, &d,
		&tx, &ty)
	m.Init(d, -b,
		-c, a,
		c*ty-d*tx, b*tx-a*ty)
}

func (matrix *Matrix) Invert() Status {
	/* Simple scaling|translation matrices are quite common... */
	if matrix.XY == 0 && matrix.YX == 0 {
		matrix.X0 = -matrix.X0
		matrix.Y0 = -matrix.Y0

		if matrix.XX != 1 {
			if matrix.XX == 0 {
				return StatusInvalidMatrix.error()
			}

			matrix.XX = 1 / matrix.XX
			matrix.X0 *= matrix.XX
		}

		if matrix.YY != 1 {
			if matrix.YY == 0 {
				return StatusInvalidMatrix.error()
			}
			matrix.YY = 1 / matrix.YY
			matrix.Y0 *= matrix.YY
		}

		return StatusSuccess
	}

	/* inv (A) = 1/det (A) * adj (A) */
	det := matrix.computeDeterminat()

	if math.IsInf(det, 0) {
		return StatusInvalidMatrix.error()
	}

	if det == 0 {
		return StatusInvalidMatrix.error()
	}

	matrix.computeAdjoint()
	matrix.scalarMultiply(1 / det)

	return StatusSuccess
}

func (m *Matrix) isInvertible() bool {
	det := m.computeDeterminat()
	// det is finite 有限的
	return !math.IsInf(det, 0) && det != 0
}

func (m *Matrix) isScale0() bool {
	return m.XX == 0 &&
		m.XY == 0 &&
		m.YX == 0 &&
		m.YY == 0
}

func (m *Matrix) computeDeterminat() float64 {
	a := m.XX
	b := m.YX
	c := m.XY
	d := m.YY
	return a*d - b*c
}

func (m *Matrix) computeBasicScaleFactors(basisScale, normalScale *float64, xBasis bool) Status {
	det := m.computeDeterminat()

	if math.IsInf(det, 0) {
		return StatusInvalidMatrix.error()
	}

	if det == 0 {
		*basisScale = 0
		*normalScale = 0
	} else {
		var x, y float64
		if xBasis {
			x = 1
		}
		if x == 0 {
			y = 1
		}
		var major, minor float64

		m.TransformDistance(&x, &y)
		major = math.Hypot(x, y)
		/*
		 * ignore mirroring
		 */
		if det < 0 {
			det = -det
		}

		if major != 0 {
			minor = det / major
		} else {
			minor = 0
		}

		if xBasis {
			*basisScale = major
			*normalScale = minor
		} else {
			*basisScale = minor
			*normalScale = major
		}
	}

	return StatusSuccess
}

func (m *Matrix) isIntegerTranslation(itx, ity *int) bool {
	if m.isTranslation() {
		x0Fixed := fixedFromDouble(m.X0)
		y0Fixed := fixedFromDouble(m.Y0)

		if x0Fixed.IsInteger() && y0Fixed.IsInteger() {
			if itx != nil {
				*itx = x0Fixed.integerPart()
			}
			if ity != nil {
				*ity = x0Fixed.integerPart()
			}
			return true
		}
	}
	return false
}

var scalingEpsilon = fixed(1).toDouble()

func (m *Matrix) hasUnityScale() bool {
	/* check that the determinant is near +/-1 */
	det := m.computeDeterminat()
	if math.Abs(det*det-1) < scalingEpsilon {
		/* check that one axis is close to zero */
		if math.Abs(m.XY) < scalingEpsilon &&
			math.Abs(m.YX) < scalingEpsilon {
			return true
		}
		if math.Abs(m.XX) < scalingEpsilon &&
			math.Abs(m.YY) < scalingEpsilon {
			return true
		}
		/* If rotations are allowed then it must instead test for
		 * orthogonality. This is xx*xy+yx*yy ~= 0.
		 */
	}
	return false
}

func (m *Matrix) isPixelExact() bool {
	if !m.hasUnityScale() {
		return false
	}

	x0Fixed := fixedFromDouble(m.X0)
	y0Fixed := fixedFromDouble(m.Y0)

	return x0Fixed.IsInteger() && y0Fixed.IsInteger()
}

func (m *Matrix) transformedCircleMajorAxis(radius float64) float64 {
	if m.hasUnityScale() {
		return radius
	}
	var a, b, c, d, f, g, h, i, j float64
	m.getAffine(&a, &b,
		&c, &d,
		nil, nil)

	i = a*a + b*b
	j = c*c + d*d

	f = 0.5 * (i + j)
	g = 0.5 * (i - j)
	h = a*c + b*d

	return radius * math.Sqrt(f+math.Hypot(g, h))
	/*
	* we don't need the minor axis length, which is
	* double min = radius * sqrt (f - sqrt (g*g+h*h));
	 */
}
