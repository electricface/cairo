package cairo

import "math"

var emptyRectangle = RectangleInt{}

var unboundedRectangle = RectangleInt{
	X:      rectIntMin,
	Y:      rectIntMin,
	Width:  rectIntMax - rectIntMin,
	Height: rectIntMax - rectIntMin,
}

func (b *box) fromDoubles(x1, y1, x2, y2 float64) {
	b.p1.x = fixedFromDouble(x1)
	b.p1.y = fixedFromDouble(y1)
	b.p2.x = fixedFromDouble(x2)
	b.p2.y = fixedFromDouble(y2)
}

func (b *box) toDoubles() (x1, y1, x2, y2 float64) {
	return b.p1.x.toDouble(),
		b.p1.y.toDouble(),
		b.p2.x.toDouble(),
		b.p2.y.toDouble()
}

func (b *box) fromRectangle(rect *RectangleInt) {
	b.p1.x = fixedFromInt(rect.X)
	b.p1.y = fixedFromInt(rect.Y)
	b.p2.x = fixedFromInt(rect.X + rect.Width)
	b.p2.y = fixedFromInt(rect.Y + rect.Height)
}

func boxesGetExtents(boxes []box, extents *box) {
	if len(boxes) == 0 {
		panic("assert failed len(boxes) > 0")
	}
	*extents = boxes[0]
	for _, b := range boxes[1:] {
		extents.addBox(&b)
	}
}

func (b *box) roundToRectangle(rectangle *RectangleInt) {
	rectangle.X = b.p1.x.integerFloor()
	rectangle.Y = b.p1.y.integerFloor()
	rectangle.Width = b.p2.x.integerCeil() - rectangle.X
	rectangle.Height = b.p2.y.integerCeil() - rectangle.Y
}

func (dst *RectangleInt) intersect(src *RectangleInt) bool {
	x1 := maxInt(dst.X, src.X)
	y1 := maxInt(dst.Y, src.Y)
	/* Beware the unsigned promotion, fortunately we have bits to spare
	 * as (CAIRO_RECT_INT_MAX - CAIRO_RECT_INT_MIN) < UINT_MAX
	 */
	x2 := minInt(dst.X+dst.Width, src.X+src.Width)
	y2 := minInt(dst.Y+dst.Height, src.Y+src.Height)

	if x1 >= x2 || y1 >= y2 {
		dst.X = 0
		dst.Y = 0
		dst.Width = 0
		dst.Height = 0

		return false
	} else {
		dst.X = x1
		dst.Y = y1
		dst.Width = x2 - x1
		dst.Height = y2 - y1

		return true
	}
}

func (dst *RectangleInt) union(src *RectangleInt) {
	x1 := minInt(dst.X, src.X)
	y1 := minInt(dst.Y, src.Y)
	/* Beware the unsigned promotion, fortunately we have bits to spare
	 * as (CAIRO_RECT_INT_MAX - CAIRO_RECT_INT_minInt) < UINT_MAX
	 */
	x2 := maxInt(dst.X+dst.Width, src.X+src.Width)
	y2 := maxInt(dst.Y+dst.Height, src.Y+src.Height)

	dst.X = x1
	dst.Y = y1
	dst.Width = x2 - x1
	dst.Height = y2 - y1
}

func (b *box) intersectsLineSegment(line *line) bool {
	var t1, t2, t3, t4 fixed

	if b.containsPoint(&line.p1) ||
		b.containsPoint(&line.p2) {
		return true
	}

	xLen := line.p2.x - line.p1.x
	yLen := line.p2.y - line.p1.y

	if xLen != 0 {
		if xLen > 0 {
			t1 = b.p1.x - line.p1.x
			t2 = b.p2.x - line.p1.x
		} else {
			t1 = line.p1.x - b.p2.x
			t2 = line.p1.x - b.p1.x
			xLen = -xLen
		}

		if (t1 < 0 || t1 > xLen) && (t2 < 0 || t2 > xLen) {
			return false
		}
	} else {
		/* Fully vertical line -- check that X is in bounds */
		if line.p1.x < b.p1.x || line.p1.x > b.p2.x {
			return false
		}
	}

	if yLen != 0 {
		if yLen > 0 {
			t3 = b.p1.y - line.p1.y
			t4 = b.p2.y - line.p1.y
		} else {
			t3 = line.p1.y - b.p2.y
			t4 = line.p1.y - b.p1.y
			yLen = -yLen
		}

		if (t3 < 0 || t3 > yLen) && (t4 < 0 || t4 > yLen) {
			return false
		}
	} else {
		/* Fully horizontal line -- check Y */
		if line.p1.y < b.p1.y || line.p1.y > b.p2.y {
			return false
		}
	}

	/* If we had a horizontal or vertical line, then it's already been checked */
	if line.p1.x == line.p2.x || line.p1.y == line.p2.y {
		return true
	}

	/* Check overlap.  Note that t1 < t2 and t3 < t4 here. */
	t1y := int64(t1) * int64(yLen)
	t2y := int64(t2) * int64(yLen)
	t3x := int64(t3) * int64(xLen)
	t4x := int64(t4) * int64(xLen)

	if t1y < t4x && t3x < t2y {
		return true
	}
	return false
}

func boxAddSplinePoint(closure interface{}, point *point, tangent *slope) Status {
	box := closure.(*box)
	box.addPoint(point)
	return StatusSuccess
}

func (extents *box) addCurveTo(a, b, c, d *point) {
	extents.addPoint(d)
	if !extents.containsPoint(b) || !extents.containsPoint(c) {
		status := splineBound(boxAddSplinePoint, extents, a, b, c, d)
		if status != StatusSuccess {
			panic("assert failed status == StatusSuccess")
		}
	}
}

func (rectI *RectangleInt) fromDouble(rectF *Rectangle) {
	rectI.X = int(math.Floor(rectF.X))
	rectI.Y = int(math.Floor(rectF.Y))
	rectI.Width = int(math.Ceil(rectF.X+rectF.Width)) - int(math.Floor(rectF.X))
	rectI.Height = int(math.Ceil(rectF.Y+rectF.Height)) - int(math.Floor(rectF.Y))
}
