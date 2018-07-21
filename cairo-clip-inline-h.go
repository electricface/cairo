package cairo

func (c *clip) isAllClipped() bool {
	return c == &clipAll
}

func (c *clip) setAllClipped() *clip {
	c.destroy()
	return &clipAll
}

func (c *clip) copyIntersectRectangle(r *RectangleInt) *clip {
	return c.copy().intersectRectangle(r)
}

func (c *clip) copyIntersectClip(other *clip) *clip {
	return c.copy().intersectClip(other)
}

func (c *clip) stealBoxes(boxes *boxes) {
	boxes.initForArray(c.boxes)
	c.boxes = nil
}

func (c *clip) unstealBoxes(boxes *boxes) {
	c.boxes = boxes.chunks.Front().Value.(boxesChunk).boxes
}
