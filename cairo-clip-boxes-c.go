package cairo

func pot(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func (c *clip) containsRectangleBox(rect *RectangleInt, box *box) bool {
	/* c == NULL means no clip, so the clip contains everything */
	if c == nil {
		return true
	}

	if c.isAllClipped() {
		return false
	}

	/* If we have a non-trivial path, just say no */
	if c.path != nil {
		return false
	}

	if !c.extents.containsRectangle(rect) {
		return false
	}

	if len(c.boxes) == 0 {
		return true
	}

	/* Check for a clip-box that wholly contains the rectangle */
	for _, clipBox := range c.boxes {
		if box.p1.x >= clipBox.p1.x &&
			box.p1.y >= clipBox.p1.y &&
			box.p2.x <= clipBox.p2.x &&
			box.p2.y <= clipBox.p2.y {
			return true
		}
	}
	return false
}

func (c *clip) containsBox(box *box) bool {
	var rect RectangleInt
	box.roundToRectangle(&rect)
	return c.containsRectangleBox(&rect, box)
}

func (c *clip) containsRectangle(rect *RectangleInt) bool {
	var box box
	box.fromRectangleInt(rect)
	return c.containsRectangleBox(rect, &box)
}

func (c *clip) intersectRectilinearPath(path *pathFixed, fillRule FillRule,
	antialias Antialias) *clip {
	var status Status
	var boxes boxes
	boxes.init()
	status = path.fillRectilinearToBoxes(fillRule, antialias, &boxes)
	if status == StatusSuccess && boxes.numBoxes != 0 {
		c = c.intersectBoxes(&boxes)
	} else {
		c = c.setAllClipped()
	}
	boxes.fini()
	return c
}

func (c *clip) intersectRectangleBox(r *RectangleInt, box0 *box) *clip {
	var extentsBox box
	changed := false

	if c == nil {
		c = clipCreate()
	}
	if c == nil {
		return c.setAllClipped()
	}

	if len(c.boxes) == 0 {
		c.boxes = c.embeddedBox[:]
		c.boxes[0] = *box0
		if c.path == nil {
			c.extents = *r
		} else {
			if !c.extents.intersect(r) {
				return c.setAllClipped()
			}
		}

		if c.path == nil {
			c.isRegion = box0.isPixelAligned()
		}
		return c
	}
	/* Does the new box wholly subsume the clip? Perform a cheap check
	 * for the common condition of a single clip rectangle.
	 */
	if len(c.boxes) == 1 &&
		c.boxes[0].p1.x >= box0.p1.x &&
		c.boxes[0].p1.y >= box0.p1.y &&
		c.boxes[0].p2.x <= box0.p2.x &&
		c.boxes[0].p2.y <= box0.p2.y {
		return c
	}

	var i, j int
	for ; i < len(c.boxes); i++ {
		b := &c.boxes[j]

		if j != i {
			*b = c.boxes[i]
		}

		if box0.p1.x > b.p1.x {
			b.p1.x = box0.p1.x
			changed = true
		}
		if box0.p2.x < b.p2.x {
			b.p2.x = box0.p2.x
			changed = true
		}

		if box0.p1.y > b.p1.y {
			b.p1.y = box0.p1.y
			changed = true
		}
		if box0.p2.y < b.p2.y {
			b.p2.y = box0.p2.y
			changed = true
		}

		if b.p2.x > b.p1.x && b.p2.y > b.p1.y {
			j++
		}
	}
	c.boxes = c.boxes[:j]

	if len(c.boxes) == 0 {
		return c.setAllClipped()
	}

	if !changed {
		return c
	}

	extentsBox = c.boxes[0]
	for _, clipBox := range c.boxes[1:] {
		if clipBox.p1.x < extentsBox.p1.x {
			extentsBox.p1.x = clipBox.p1.x
		}

		if clipBox.p1.y < extentsBox.p1.y {
			extentsBox.p1.y = clipBox.p1.y
		}

		if clipBox.p2.x > extentsBox.p2.x {
			extentsBox.p2.x = clipBox.p2.x
		}

		if clipBox.p2.y > extentsBox.p2.y {
			extentsBox.p2.y = clipBox.p2.y
		}
	}

	if c.path == nil {
		extentsBox.roundToRectangle(&c.extents)
	} else {
		var extentsRect RectangleInt
		extentsBox.roundToRectangle(&extentsRect)
		if !c.extents.intersect(&extentsRect) {
			return c.setAllClipped()
		}
	}

	if c.region != nil {
		c.region.destroy()
		c.region = nil
	}

	c.isRegion = false
	return c
}

func (c *clip) intersectBox(box *box) *clip {
	if c.isAllClipped() {
		return c
	}

	var r RectangleInt
	box.roundToRectangle(&r)
	if r.Width == 0 || r.Height == 0 {
		return c.setAllClipped()
	}

	return c.intersectRectangleBox(&r, box)
}

func (boxes *boxes) copyToClip(clip *clip) bool {
	/* XXX cow-boxes? */
	if boxes.numBoxes == 1 {
		clip.boxes = clip.embeddedBox[:]
		clip.boxes[0] = boxes.chunks.Front().Value.(boxesChunk).boxes[0]
		return true
	}

	clip.boxes = boxes.toArray()
	if clip.boxes == nil {
		clip.setAllClipped()
		return false
	}
	return true
}

func (c *clip) intersectBoxes(boxes0 *boxes) *clip {
	var clipBoxes boxes
	var limits box
	var extents RectangleInt

	if c.isAllClipped() {
		return c
	}

	if boxes0.numBoxes == 0 {
		return c.setAllClipped()
	}

	if boxes0.numBoxes == 1 {
		firstBox := &boxes0.chunks.Front().Value.(boxesChunk).boxes[0]
		return c.intersectBox(firstBox)
	}

	if c == nil {
		c = clipCreate()
	}

	if len(c.boxes) != 0 {
		clipBoxes.initForArray(c.boxes)
		if clipBoxes.intersect(boxes0, &clipBoxes) {
			c = c.setAllClipped()
			goto out
		}

		c.boxes = nil
		boxes0 = &clipBoxes
	}

	if boxes0.numBoxes == 0 {
		c = c.setAllClipped()
		goto out
	}

	boxes0.copyToClip(c)
	boxes0.extents(&limits)
	limits.roundToRectangle(&extents)
	if c.path == nil {
		c.extents = extents
	} else if !c.extents.intersect(&extents) {
		c = c.setAllClipped()
		goto out
	}

	if c.region != nil {
		c.region.destroy()
		c.region = nil
	}
	c.isRegion = false

out:
	if boxes0 == &clipBoxes {
		clipBoxes.fini()
	}
	return c
}

func (c *clip) intersectRectangle(r *RectangleInt) *clip {
	if c.isAllClipped() {
		return c
	}

	if r.Width == 0 || r.Height == 0 {
		return c.setAllClipped()
	}

	var box box
	box.fromRectangleInt(r)
	return c.intersectRectangleBox(r, &box)
}

type reduce struct {
	clip         *clip
	limit        box
	extents      box
	inside       bool
	currentPoint point
	lastMoveTo   point
}

func addClippedEdge(r *reduce, p1, p2 *point, y1, y2 int) {
	var x fixed
	x = edgeComputeIntersectionXForY(p1, p2, fixed(y1))
	if x < r.extents.p1.x {
		r.extents.p1.x = x
	}

	x = edgeComputeIntersectionXForY(p1, p2, fixed(y2))
	if x > r.extents.p2.x {
		r.extents.p2.x = x
	}

	if fixed(y1) < r.extents.p1.y {
		r.extents.p1.y = fixed(y1)
	}

	if fixed(y2) > r.extents.p2.y {
		r.extents.p2.y = fixed(y2)
	}
	r.inside = true
}

func addEdge(r *reduce, p1, p2 *point) {
	var top, bottom int
	var topY, botY int
	var n int

	if p1.y < p2.y {
		top = int(p1.y)
		bottom = int(p2.y)
	} else {
		top = int(p2.y)
		bottom = int(p1.y)
	}

	if bottom < int(r.limit.p1.y) || top > int(r.limit.p2.y) {
		return
	}

	if p1.x > p2.x {
		p1, p2 = p2, p1
	}

	if p2.x <= r.limit.p1.x || p1.x >= r.limit.p2.x {
		return
	}

	for n = 0; n < len(r.clip.boxes); n++ {
		limits := &r.clip.boxes[n]
		if bottom < int(limits.p1.y) || top > int(limits.p2.y) {
			continue
		}

		if p2.x <= limits.p1.x || p1.x >= limits.p2.x {
			continue
		}

		if p1.x >= limits.p1.x && p2.x <= limits.p1.x {
			topY = top
			botY = bottom
		} else {
			var p1Y, p2Y int
			p1Y = int(edgeComputeIntersectionYForX(p1, p2, limits.p1.x))

			p2Y = int(edgeComputeIntersectionYForX(p1, p2, limits.p2.x))

			if p1Y < p2Y {
				topY = p1Y
				botY = p2Y
			} else {
				topY = p2Y
				botY = p1Y
			}

			if topY < top {
				topY = top
			}
			if botY > bottom {
				botY = bottom
			}
		}

		if topY < int(limits.p1.y) {
			topY = int(limits.p1.y)
		}

		if botY > int(limits.p2.y) {
			botY = int(limits.p2.y)
		}
		if botY > topY {
			addClippedEdge(r, p1, p2, topY, botY)
		}
	}

}

func reduceLineTo(closure interface{}, point *point) Status {
	r := closure.(*reduce)
	addEdge(r, &r.currentPoint, point)
	r.currentPoint = *point
	return StatusSuccess
}

func reduceClose(closure interface{}) Status {
	r := closure.(*reduce)
	return reduceLineTo(r, &r.lastMoveTo)
}

func reduceMoveTo(closure interface{}, point *point) Status {
	r := closure.(*reduce)

	/* close current subpath */
	status := reduceClose(closure)

	/* make sure that the closure represents a degenerate path */
	r.currentPoint = *point
	r.lastMoveTo = *point

	return status
}

func (c *clip) reduceToBoxes() *clip {
	//var r reduce
	//var clipPath *clipPath
	//var status Status

	return c
	// TODO: 这段代码很奇特，直接返回了 c
	// 我在 gitlab.com 上提交了一个 issues https://gitlab.com/cairo/cairo/issues/2
}

func (c *clip) reduceToRectangle(r *RectangleInt) *clip {
	if c.isAllClipped() {
		return c
	}

	if c.containsRectangle(r) {
		return (*clip)(nil).intersectRectangle(r)
	}

	copy0 := c.copyIntersectRectangle(r)
	if copy0.isAllClipped() {
		return copy0
	}

	return copy0.reduceToBoxes()
}

func (c *clip) reduceForComposite(extents *compoisteRectangles) *clip {
	var r *RectangleInt
	if extents.isBounded {
		r = &extents.bounded
	} else {
		r = &extents.unbounded
	}
	return c.reduceToRectangle(r)
}

func clipFromBoxes(boxes *boxes) *clip {
	c := clipCreate()
	if c == nil {
		return c.setAllClipped()
	}

	if !boxes.copyToClip(c) {
		return c
	}

	var extents box
	boxes.extents(&extents)
	extents.roundToRectangle(&c.extents)
	return c
}
