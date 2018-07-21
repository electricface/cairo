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
	box.fromRectangle(rect)
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
