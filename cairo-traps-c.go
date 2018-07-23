package cairo

func (ts *traps) init() {
	ts.status = StatusSuccess
	ts.maybeRegion = true
	ts.isRectilinear = false
	ts.isRectangular = false
	ts.traps = ts.trapsEmbedded[:0]
	ts.limits = nil
	ts.hasIntersections = false
}

func (ts *traps) limit(limits []box) {
	ts.limits = limits
	ts.bounds = limits[0]
	for i := 1; i < len(limits); i++ {
		ts.bounds.addBox(&limits[i])
	}
}

func (ts *traps) initWithClip(clip *clip) {
	ts.init()
	if clip != nil {
		ts.limit(clip.boxes)
	}
}

func (ts *traps) clear() {
	ts.status = StatusSuccess
	ts.maybeRegion = true
	ts.isRectilinear = false
	ts.isRectangular = false
	ts.traps = nil
	ts.hasIntersections = false
}

func (ts *traps) fini() {
}

func (ts *traps) grow() bool {
	newSize := 4 * cap(ts.traps)

	newTraps := make([]trapzoid, newSize)
	copy(newTraps, ts.traps)
	ts.traps = newTraps[:len(ts.traps)]
	return true
}

func (ts *traps) addTrap(top, bottom fixed, left, right *line) {
	if left.p1.y == left.p2.y {
		panic("assert failed left.p1.y != left.p2.y")
	}
	if right.p1.y == right.p2.y {
		panic("assert failed right.p1.y != right.p2.y")
	}
	if !(bottom > top) {
		panic("assert failed bottom > top")
	}

	if len(ts.traps) == cap(ts.traps) {
		ts.grow()
	}
	trap := trapzoid{
		top:    top,
		bottom: bottom,
		left:   *left,
		right:  *right,
	}
	ts.traps = append(ts.traps, trap)
}

func (ts *traps) addClippedTrap(_top, _bottom fixed, _left, _right *line) {
	/* Note: With the goofy trapezoid specification, (where an
	 * arbitrary two points on the lines can specified for the left
	 * and right edges), these limit checks would not work in
	 * general. For example, one can imagine a trapezoid entirely
	 * within the limits, but with two points used to specify the left
	 * edge entirely to the right of the limits.  Fortunately, for our
	 * purposes, cairo will never generate such a crazy
	 * trapezoid. Instead, cairo always uses for its points the
	 * extreme positions of the edge that are visible on at least some
	 * trapezoid. With this constraint, it's impossible for both
	 * points to be outside the limits while the relevant edge is
	 * entirely inside the limits.
	 */
	if len(ts.limits) != 0 {
		b := &ts.bounds
		top := _top
		bottom := _bottom
		left := *_left
		right := *_right

		/* Trivially reject if trapezoid is entirely to the right or
		 * to the left of the limits. */
		if left.p1.x >= b.p2.x && left.p2.x >= b.p2.x {
			return
		}

		if right.p1.x <= b.p1.x && right.p2.x <= b.p1.x {
			return
		}

		/* And reject if the trapezoid is entirely above or below */
		if top >= b.p2.y || bottom <= b.p1.y {
			return
		}

		/* Otherwise, clip the trapezoid to the limits. We only clip
		 * where an edge is entirely outside the limits. If we wanted
		 * to be more clever, we could handle cases where a trapezoid
		 * edge intersects the edge of the limits, but that would
		 * require slicing this trapezoid into multiple trapezoids,
		 * and I'm not sure the effort would be worth it. */
		if top < b.p1.y {
			top = b.p1.y
		}

		if bottom > b.p2.y {
			bottom = b.p2.y
		}

		if left.p1.x <= b.p1.x && left.p2.x <= b.p1.x {
			left.p1.x = b.p1.x
			left.p2.x = b.p1.x
		}

		if right.p1.x >= b.p2.x && right.p2.x >= b.p2.x {
			right.p1.x = b.p2.x
			right.p2.x = b.p2.x
		}
		/* Trivial discards for empty trapezoids that are likely to
		 * be produced by our tessellators (most notably convex_quad
		 * when given a simple rectangle).
		 */
		if top >= bottom {
			return
		}
		/* cheap colinearity check */
		if right.p1.x <= left.p1.x && right.p1.y == left.p1.y &&
			right.p2.x <= left.p2.x && right.p2.y == left.p2.y {
			return
		}

		ts.addTrap(top, bottom, &left, &right)
	} else {
		ts.addTrap(_top, _bottom, _left, _right)
	}
}

func comparePointFixedByY(a, b *point) int {
	var ret int = int(a.y - b.y)
	if ret == 0 {
		ret = int(a.x - b.x)
	}
	return ret
}

func (ts *traps) tessellateConvexQuad(q [4]point) {
	var a, b, c, d int
	var ab, ad slope
	var b_left_of_d bool
	var left, right line

	/* Choose a as a point with minimal y */
	a = 0
	for i := 1; i < 4; i++ {
		if comparePointFixedByY(&q[i], &q[a]) < 0 {
			a = i
		}
	}

	/* b and d are adjacent to a, while c is opposite */
	b = (a + 1) % 4
	c = (a + 2) % 4
	d = (a + 3) % 4

	/* Choose between b and d so that b.y is less than d.y */
	if comparePointFixedByY(&q[d], &q[b]) < 0 {
		b = (a + 3) % 4
		d = (a + 1) % 4
	}

	/* Without freedom left to choose anything else, we have four
	 * cases to tessellate.
	 *
	 * First, we have to determine the Y-axis sort of the four
	 * vertices, (either abcd or abdc). After that we need to detemine
	 * which edges will be "left" and which will be "right" in the
	 * resulting trapezoids. This can be determined by computing a
	 * slope comparison of ab and ad to determine if b is left of d or
	 * not.
	 *
	 * Note that "left of" here is in the sense of which edges should
	 * be the left vs. right edges of the trapezoid. In particular, b
	 * left of d does *not* mean that b.x is less than d.x.
	 *
	 * This should hopefully be made clear in the lame ASCII art
	 * below. Since the same slope comparison is used in all cases, we
	 * compute it before testing for the Y-value sort. */

	/* Note: If a == b then the ab slope doesn't give us any
	 * information. In that case, we can replace it with the ac (or
	 * equivalenly the bc) slope which gives us exactly the same
	 * information we need. At worst the names of the identifiers ab
	 * and b_left_of_d are inaccurate in this case, (would be ac, and
	 * c_left_of_d). */
	if q[a].x == q[b].x && q[a].y == q[b].y {
		ab.init(&q[a], &q[c])
	} else {
		ab.init(&q[a], &q[b])
	}

	ad.init(&q[a], &q[d])

	b_left_of_d = ab.compare(&ad) > 0

	if q[c].y <= q[d].y {
		if b_left_of_d {
			/* Y-sort is abcd and b is left of d, (slope(ab) > slope (ad))
			 *
			 *                      top bot left right
			 *        _a  a  a
			 *      / /  /|  |\      a.y b.y  ab   ad
			 *     b /  b |  b \
			 *    / /   | |   \ \    b.y c.y  bc   ad
			 *   c /    c |    c \
			 *  | /      \|     \ \  c.y d.y  cd   ad
			 *  d         d       d
			 */
			left.p1 = q[a]
			left.p2 = q[b]
			right.p1 = q[a]
			right.p2 = q[d]
			ts.addClippedTrap(q[a].y, q[b].y, &left, &right)
			left.p1 = q[b]
			left.p2 = q[c]
			ts.addClippedTrap(q[b].y, q[c].y, &left, &right)
			left.p1 = q[c]
			left.p2 = q[d]
			ts.addClippedTrap(q[c].y, q[d].y, &left, &right)
		} else {
			/* Y-sort is abcd and b is right of d, (slope(ab) <= slope (ad))
			 *
			 *       a  a  a_
			 *      /|  |\  \ \     a.y b.y  ad  ab
			 *     / b  | b  \ b
			 *    / /   | |   \ \   b.y c.y  ad  bc
			 *   / c    | c    \ c
			 *  / /     |/      \ | c.y d.y  ad  cd
			 *  d       d         d
			 */
			left.p1 = q[a]
			left.p2 = q[d]
			right.p1 = q[a]
			right.p2 = q[b]
			ts.addClippedTrap(q[a].y, q[b].y, &left, &right)
			right.p1 = q[b]
			right.p2 = q[c]
			ts.addClippedTrap(q[b].y, q[c].y, &left, &right)
			right.p1 = q[c]
			right.p2 = q[d]
			ts.addClippedTrap(q[c].y, q[d].y, &left, &right)
		}
	} else {
		if b_left_of_d {
			/* Y-sort is abdc and b is left of d, (slope (ab) > slope (ad))
			 *
			 *        a   a     a
			 *       //  / \    |\     a.y b.y  ab  ad
			 *     /b/  b   \   b \
			 *    / /    \   \   \ \   b.y d.y  bc  ad
			 *   /d/      \   d   \ d
			 *  //         \ /     \|  d.y c.y  bc  dc
			 *  c           c       c
			 */
			left.p1 = q[a]
			left.p2 = q[b]
			right.p1 = q[a]
			right.p2 = q[d]
			ts.addClippedTrap(q[a].y, q[b].y, &left, &right)
			left.p1 = q[b]
			left.p2 = q[c]
			ts.addClippedTrap(q[b].y, q[d].y, &left, &right)
			right.p1 = q[d]
			right.p2 = q[c]
			ts.addClippedTrap(q[d].y, q[c].y, &left, &right)
		} else {
			/* Y-sort is abdc and b is right of d, (slope (ab) <= slope (ad))
			 *
			 *      a     a   a
			 *     /|    / \  \\       a.y b.y  ad  ab
			 *    / b   /   b  \b\
			 *   / /   /   /    \ \    b.y d.y  ad  bc
			 *  d /   d   /	 \d\
			 *  |/     \ /         \\  d.y c.y  dc  bc
			 *  c       c	   c
			 */
			left.p1 = q[a]
			left.p2 = q[d]
			right.p1 = q[a]
			right.p2 = q[b]
			ts.addClippedTrap(q[a].y, q[b].y, &left, &right)
			right.p1 = q[b]
			right.p2 = q[c]
			ts.addClippedTrap(q[b].y, q[d].y, &left, &right)
			left.p1 = q[d]
			left.p2 = q[c]
			ts.addClippedTrap(q[d].y, q[c].y, &left, &right)
		}
	}
}

func (ts *traps) addTri(y1, y2 int, left, right *line) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	if linesCompareAtY(left, right, y1) > 0 {
		left, right = right, left
	}

	ts.addClippedTrap(fixed(y1), fixed(y2), left, right)
}

func (ts *traps) tessellateTriangleWithEdges(t [3]point, edges [4]point) {
	var lines [3]line

	if edges[0].y <= edges[1].y {
		lines[0].p1 = edges[0]
		lines[0].p2 = edges[1]
	} else {
		lines[0].p1 = edges[1]
		lines[0].p2 = edges[0]
	}

	if edges[2].y <= edges[3].y {
		lines[1].p1 = edges[2]
		lines[1].p2 = edges[3]
	} else {
		lines[1].p1 = edges[3]
		lines[1].p2 = edges[2]
	}

	if t[1].y == t[2].y {
		ts.addTri(int(t[0].y), int(t[1].y), &lines[0], &lines[1])
		return
	}

	if t[1].y <= t[2].y {
		lines[2].p1 = t[1]
		lines[2].p2 = t[2]
	} else {
		lines[2].p1 = t[2]
		lines[2].p2 = t[1]
	}

	if ((t[1].y - t[0].y) < 0) != ((t[2].y - t[0].y) < 0) {
		ts.addTri(int(t[0].y), int(t[1].y), &lines[0], &lines[2])
		ts.addTri(int(t[0].y), int(t[2].y), &lines[1], &lines[2])
	} else if absFixed(t[1].y-t[0].y) < absFixed(t[2].y-t[0].y) {
		ts.addTri(int(t[0].y), int(t[1].y), &lines[0], &lines[1])
		ts.addTri(int(t[1].y), int(t[2].y), &lines[2], &lines[1])
	} else {
		ts.addTri(int(t[0].y), int(t[2].y), &lines[1], &lines[0])
		ts.addTri(int(t[1].y), int(t[2].y), &lines[2], &lines[0])
	}
}

func absFixed(f fixed) fixed {
	if f < 0 {
		return -f
	}
	return f
}

func (ts *traps) initBoxes(boxes *boxes) Status {
	ts.init()
	ts.traps = make([]trapzoid, boxes.numBoxes)

	ts.isRectilinear = true
	ts.isRectangular = true
	ts.maybeRegion = boxes.isPixelAligned

	var trapsIdx int
	for elem := boxes.chunks.Front(); elem != nil; elem = elem.Next() {
		chunk := elem.Value.(boxesChunk)
		for _, box := range chunk.boxes {
			trap := &ts.traps[trapsIdx]

			trap.top = box.p1.y
			trap.bottom = box.p2.y

			trap.left.p1 = box.p1
			trap.left.p2.x = box.p1.x
			trap.left.p2.y = box.p2.y

			trap.right.p1.x = box.p2.x
			trap.right.p1.y = box.p1.y
			trap.right.p2 = box.p2
			trapsIdx++
		}
	}
	return StatusSuccess
}

func (ts *traps) tessellateRectangle(top_left, bottom_right *point) Status {
	var left, right line
	var top, bottom fixed

	if top_left.y == bottom_right.y {
		return StatusSuccess
	}

	if top_left.x == bottom_right.x {
		return StatusSuccess
	}

	left.p1.x = top_left.x
	left.p2.x = top_left.x

	left.p1.y = top_left.y
	right.p1.y = top_left.y

	right.p1.x = bottom_right.x
	right.p2.x = bottom_right.x

	left.p2.y = bottom_right.y
	right.p2.y = bottom_right.y

	top = top_left.y
	bottom = bottom_right.y

	if len(ts.limits) != 0 {
		if top >= ts.bounds.p2.y || bottom <= ts.bounds.p1.y {
			return StatusSuccess
		}

		/* support counter-clockwise winding for rectangular tessellation */
		reversed := top_left.x > bottom_right.x
		if reversed {
			right.p1.x = top_left.x
			right.p2.x = top_left.x

			left.p1.x = bottom_right.x
			left.p2.x = bottom_right.x
		}

		if left.p1.x >= ts.bounds.p2.x || right.p1.x <= ts.bounds.p1.x {
			return StatusSuccess
		}

		for n := 0; n < len(ts.limits); n++ {
			limits := &ts.limits[n]
			var _left, _right line
			var _top, _bottom fixed

			if top >= limits.p2.y {
				continue
			}
			if bottom <= limits.p1.y {
				continue
			}

			/* Trivially reject if trapezoid is entirely to the right or
			 * to the left of the limits. */
			if left.p1.x >= limits.p2.x {
				continue
			}
			if right.p1.x <= limits.p1.x {
				continue
			}

			/* Otherwise, clip the trapezoid to the limits. */
			_top = top
			if _top < limits.p1.y {
				_top = limits.p1.y
			}

			_bottom = bottom
			if _bottom > limits.p2.y {
				_bottom = limits.p2.y
			}

			if _bottom <= _top {
				continue
			}

			_left = left
			if _left.p1.x < limits.p1.x {
				_left.p1.x = limits.p1.x
				_left.p1.y = limits.p1.y
				_left.p2.x = limits.p1.x
				_left.p2.y = limits.p2.y
			}

			_right = right
			if _right.p1.x > limits.p2.x {
				_right.p1.x = limits.p2.x
				_right.p1.y = limits.p1.y
				_right.p2.x = limits.p2.x
				_right.p2.y = limits.p2.y
			}

			if left.p1.x >= right.p1.x {
				continue
			}

			if reversed {
				ts.addTrap(_top, _bottom, &_right, &_left)
			} else {
				ts.addTrap(_top, _bottom, &_left, &_right)
			}
		}
	} else {
		ts.addTrap(top, bottom, &left, &right)
	}

	return ts.status
}
