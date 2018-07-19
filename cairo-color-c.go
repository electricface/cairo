package cairo

var colorWhite = color{
	1, 1, 1, 1,
	0xffff, 0xffff, 0xffff, 0xffff,
}

var colorBlack = color{
	0, 0, 0, 1,
	0, 0, 0, 0xffff,
}

var colorTransparent = color{}

var colorMagenta = color{
	1, 0, 1, 1,
	0xffff, 0, 0xffff, 0xffff,
}

func stockColor(stock stock) *color {
	switch stock {
	case stockWhite:
		return &colorWhite
	case stockBlack:
		return &colorBlack
	case stockTransparent:
		return &colorTransparent
	default:
		assertNotReached()
		return &colorMagenta
	}
}

func colorDoubleToShort(d float64) uint16 {
	return uint16(d*65535.0 + 0.5)
}

func (c *color) computeShorts() {
	c.redShort = colorDoubleToShort(c.red * c.alpha)
	c.greenShort = colorDoubleToShort(c.green * c.alpha)
	c.blueShort = colorDoubleToShort(c.blue * c.alpha)
	c.alphaShort = colorDoubleToShort(c.alpha)
}

func (c *color) initRGBA(red, green, blue, alpha float64) {
	c.red = red
	c.green = green
	c.blue = blue
	c.alpha = alpha
	c.computeShorts()
}

func (c *color) multiplyAlpha(alpha float64) {
	c.alpha = alpha
	c.computeShorts()
}

func (c *color) getRGBA() (red, green, blue, alpha float64) {
	return c.red, c.green, c.blue, c.alpha
}

func (c *color) getRGBAPremultiplied() (red, green, blue, alpha float64) {
	return c.red * c.alpha,
		c.green * c.alpha,
		c.blue * c.alpha,
		c.alpha
}

func (a *color) equal(b *color) bool {
	if a == b {
		return true
	}

	if a.alphaShort != b.alphaShort {
		return false
	}

	if a.alphaShort == 0 {
		return true
	}

	return a.redShort == b.redShort &&
		a.greenShort == b.greenShort &&
		a.blueShort == b.blueShort
}

func (a *colorStop) equal(b *colorStop) bool {
	if a == b {
		return true
	}

	return a.alphaShort == b.alphaShort &&
		a.redShort == b.redShort &&
		a.greenShort == b.greenShort &&
		a.blueShort == b.blueShort
}

func (c *color) getContent() Content {
	if c.isOpaque() {
		return ContentColor
	}

	if c.redShort == 0 &&
		c.greenShort == 0 &&
		c.blueShort == 0 {
		return ContentAlpha
	}

	return ContentColorAlpha
}
