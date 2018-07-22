package cairo

import "math"

func splineIntersects(a, b, c, d *point, box0 box) bool {
	var bounds box
	if box0.containsPoint(a) ||
		box0.containsPoint(b) ||
		box0.containsPoint(c) ||
		box0.containsPoint(d) {
		return true
	}
	bounds.p2 = *a
	bounds.p1 = *a
	bounds.addPoint(b)
	bounds.addPoint(c)
	bounds.addPoint(d)

	if bounds.p2.x <= box0.p1.x || bounds.p1.x >= box0.p2.x ||
		bounds.p2.y <= box0.p1.y || bounds.p1.y >= box0.p2.y {
		return false
	}

	return true
}

func (s *spline) init(addPointFunc splineAddPointFunc, closure interface{}, a, b, c, d *point) bool {
	/* If both tangents are zero, this is just a straight line */
	if a.x == b.x && a.y == b.y && c.x == d.x && c.y == d.y {
		return false
	}

	s.addPointFunc = addPointFunc
	s.closure = closure

	s.knots.a = *a
	s.knots.b = *b
	s.knots.c = *c
	s.knots.d = *d

	if a.x != b.x || a.y != b.y {
		s.initialSlope.init(&s.knots.a, &s.knots.b)
	} else if a.x != c.x || a.y != c.y {
		s.initialSlope.init(&s.knots.a, &s.knots.c)
	} else if a.x != d.x || a.y != d.y {
		s.initialSlope.init(&s.knots.a, &s.knots.d)
	} else {
		return false
	}

	if c.x != d.x || c.y != d.y {
		s.finalSlope.init(&s.knots.c, &s.knots.d)
	} else if b.x != d.x || b.y != d.y {
		s.finalSlope.init(&s.knots.b, &s.knots.d)
	} else {
		return false /* just treat this as a straight-line from a . d */
	}

	/* XXX if the initial, final and vector are all equal, this is just a line */
	return true
}

func (s *spline) addPoint(point0, knot *point) Status {
	var prev *point
	var slope0 slope

	prev = &s.lastPoint

	if prev.x == point0.x && prev.y == point0.y {
		return StatusSuccess
	}

	slope0.init(point0, knot)
	s.lastPoint = *point0

	return s.addPointFunc(point0, &slope0)
}

func lerpHalf(a, b, result *point) {
	result.x = a.x + ((b.x - a.x) >> 1)
	result.y = a.y + ((b.y - a.y) >> 1)
}

func deCasteljau(s1, s2 *splineKnots) {
	var ab, bc, cd, abbc, bccd, final point

	lerpHalf(&s1.a, &s1.b, &ab)
	lerpHalf(&s1.b, &s1.c, &bc)
	lerpHalf(&s1.c, &s1.d, &cd)
	lerpHalf(&ab, &bc, &abbc)
	lerpHalf(&bc, &cd, &bccd)
	lerpHalf(&abbc, &bccd, &final)

	s2.a = final
	s2.b = bccd
	s2.c = cd
	s2.d = s1.d

	s1.b = ab
	s1.c = abbc
	s1.d = final
}

func splineErrorSquared(knots *splineKnots) float64 {
	var bdx, bdy, berr float64
	var cdx, cdy, cerr float64

	/* We are going to compute the distance (squared) between each of the the b
	 * and c control points and the segment a-b. The maximum of these two
	 * distances will be our approximation error. */

	bdx = (knots.b.x - knots.a.x).toDouble()
	bdy = (knots.b.y - knots.a.y).toDouble()

	cdx = (knots.c.x - knots.a.x).toDouble()
	cdy = (knots.c.y - knots.a.y).toDouble()

	if knots.a.x != knots.d.x || knots.a.y != knots.d.y {
		/* Intersection point (px):
		 *     px = p1 + u(p2 - p1)
		 *     (p - px) ∙ (p2 - p1) = 0
		 * Thus:
		 *     u = ((p - p1) ∙ (p2 - p1)) / ∥p2 - p1∥²
		 */

		var dx, dy, u, v float64

		dx = (knots.d.x - knots.a.x).toDouble()
		dy = (knots.d.y - knots.a.y).toDouble()
		v = dx*dx + dy*dy

		u = bdx*dx + bdy*dy
		if u <= 0 {
			/* bdx -= 0
			 * bdy -= 0
			 */
		} else if u >= v {
			bdx -= dx
			bdy -= dy
		} else {
			bdx -= u / v * dx
			bdy -= u / v * dy
		}

		u = cdx*dx + cdy*dy
		if u <= 0 {
			/* cdx -= 0
			 * cdy -= 0
			 */
		} else if u >= v {
			cdx -= dx
			cdy -= dy
		} else {
			cdx -= u / v * dx
			cdy -= u / v * dy
		}
	}

	berr = bdx*bdx + bdy*bdy
	cerr = cdx*cdx + cdy*cdy
	if berr > cerr {
		return berr
	}
	return cerr
}

func splineDecomposeInto(s1 *splineKnots, toleranceSquared float64, result *spline) Status {
	if splineErrorSquared(s1) < toleranceSquared {
		return result.addPoint(&s1.a, &s1.b)
	}

	var s2 splineKnots
	deCasteljau(s1, &s2)

	status := splineDecomposeInto(s1, toleranceSquared, result)
	if status != 0 {
		return status
	}

	return splineDecomposeInto(&s2, toleranceSquared, result)
}

func (s *spline) decompose(tolerance float64) Status {
	s1 := s.knots
	s.lastPoint = s1.a
	status := splineDecomposeInto(&s1, tolerance*tolerance, s)
	if status != 0 {
		return status
	}

	return s.addPointFunc(s.closure, &s.knots.d, &s.finalSlope)
}

/* Note: this function is only good for computing bounds in device space. */
func splineBound(addPointFunc splineAddPointFunc, closure interface{}, p0, p1, p2, p3 *point) Status {
	var x0, x1, x2, x3 float64
	var y0, y1, y2, y3 float64
	var a, b, c float64
	var t [4]float64
	var tNum int
	var i int

	var status Status

	x0 = (p0.x).toDouble()
	y0 = (p0.y).toDouble()
	x1 = (p1.x).toDouble()
	y1 = (p1.y).toDouble()
	x2 = (p2.x).toDouble()
	y2 = (p2.y).toDouble()
	x3 = (p3.x).toDouble()
	y3 = (p3.y).toDouble()

	/* The spline can be written as a polynomial of the four points:
	 *
	 *   (1-t)³p0 + 3t(1-t)²p1 + 3t²(1-t)p2 + t³p3
	 *
	 * for 0≤t≤1.  Now, the X and Y components of the spline follow the
	 * same polynomial but with x and y replaced for p.  To find the
	 * bounds of the spline, we just need to find the X and Y bounds.
	 * To find the bound, we take the derivative and equal it to zero,
	 * and solve to find the t's that give the extreme points.
	 *
	 * Here is the derivative of the curve, sorted on t:
	 *
	 *   3t²(-p0+3p1-3p2+p3) + 2t(3p0-6p1+3p2) -3p0+3p1
	 *
	 * Let:
	 *
	 *   a = -p0+3p1-3p2+p3
	 *   b =  p0-2p1+p2
	 *   c = -p0+p1
	 *
	 * Gives:
	 *
	 *   a.t² + 2b.t + c = 0
	 *
	 * With:
	 *
	 *   delta = b*b - a*c
	 *
	 * the extreme points are at -c/2b if a is zero, at (-b±√delta)/a if
	 * delta is positive, and at -b/a if delta is zero.
	 */
	add := func(t0 float64) {
		_t0 := t0
		if 0 < _t0 && _t0 < 1 {
			t[tNum] = _t0
			tNum++
		}
	}

	findExtremes := func(a, b, c float64) {
		if a == 0 {
			if b != 0 {
				add(-c / (2 * b))
			}
		} else {
			b2 := b * b
			delta := b2 - a*c
			if delta > 0 {
				var feasible bool
				_2ab := 2 * a * b
				/* We are only interested in solutions t that satisfy 0<t<1
				 * here.  We do some checks to avoid sqrt if the solutions
				 * are not in that range.  The checks can be derived from:
				 *
				 *   0 < (-b±√delta)/a < 1
				 */

				if _2ab >= 0 {
					feasible = delta > b2 && delta < a*a+b2+_2ab
				} else if -b/a >= 1 {
					feasible = delta < b2 && delta > a*a+b2+_2ab
				} else {
					feasible = delta < b2 || delta < a*a+b2+_2ab
				}

				if feasible {
					sqrtDelta := math.Sqrt(delta)
					add((-b - sqrtDelta) / a)
					add((-b + sqrtDelta) / a)
				}
			} else if delta == 0 {
				add(-b / a)
			}
		}
	}
	/* Find X extremes */
	a = -x0 + 3*x1 - 3*x2 + x3
	b = x0 - 2*x1 + x2
	c = -x0 + x1
	findExtremes(a, b, c)

	/* Find Y extremes */
	a = -y0 + 3*y1 - 3*y2 + y3
	b = y0 - 2*y1 + y2
	c = -y0 + y1
	findExtremes(a, b, c)

	status = addPointFunc(closure, p0, nil)
	if status != 0 {
		return status
	}

	for i = 0; i < tNum; i++ {
		var p point
		var x, y float64

		var t_1_0, t_0_1 float64
		var t_2_0, t_0_2 float64
		var t_3_0, t_2_1_3, t_1_2_3, t_0_3 float64

		t_1_0 = t[i]      /*      t  */
		t_0_1 = 1 - t_1_0 /* (1 - t) */

		t_2_0 = t_1_0 * t_1_0 /*      t  *      t  */
		t_0_2 = t_0_1 * t_0_1 /* (1 - t) * (1 - t) */

		t_3_0 = t_2_0 * t_1_0       /*      t  *      t  *      t      */
		t_2_1_3 = t_2_0 * t_0_1 * 3 /*      t  *      t  * (1 - t) * 3 */
		t_1_2_3 = t_1_0 * t_0_2 * 3 /*      t  * (1 - t) * (1 - t) * 3 */
		t_0_3 = t_0_1 * t_0_2       /* (1 - t) * (1 - t) * (1 - t)     */

		/* Bezier polynomial */
		x = x0*t_0_3 +
			x1*t_1_2_3 +
			x2*t_2_1_3 +
			x3*t_3_0

		y = y0*t_0_3 +
			y1*t_1_2_3 +
			y2*t_2_1_3 +
			y3*t_3_0

		p.x = fixedFromDouble(x)
		p.y = fixedFromDouble(y)
		status = addPointFunc(closure, &p, nil)
		if status == 0 {
			return status
		}
	}

	return addPointFunc(closure, p3, nil)
}
