package cairo

func (c *clip) extractRegion() {
	if len(c.boxes) == 0 {
		return
	}

	var r []RectangleInt
	r = make([]RectangleInt, len(c.boxes))

	isRegion := c.path == nil
	var i int
	for i = 0; i < len(c.boxes); i++ {
		b := &c.boxes[i]
		if isRegion {
			isRegion = (b.p1.x | b.p1.y | b.p2.x | b.p2.y).IsInteger()
		}
		r[i].X = b.p1.x.integerFloor()
		r[i].Y = b.p1.y.integerFloor()
		r[i].Width = b.p2.x.integerCeil() - r[i].X
		r[i].Height = b.p2.y.integerCeil() - r[i].Y
	}
	c.isRegion = isRegion
	c.region = regionCreateRectangles(r, i)
}

func (c *clip) getRegion() *region {
	if c == nil {
		return nil
	}

	if c.region == nil {
		c.extractRegion()
	}
	return c.region
}

func (c *clip) getIsRegion() bool {
	if c == nil {
		return true
	}

	if c.isRegion {
		return true
	}

	/* XXX Geometric reduction? */

	if c.path != nil {
		return false
	}

	if len(c.boxes) == 0 {
		return true
	}

	if c.region == nil {
		c.extractRegion()
	}

	return c.isRegion
}
