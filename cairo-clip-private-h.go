package cairo

type clipPath struct {
	path      pathFixed
	fillRule  FillRule
	tolerance float64
	antialias Antialias
	prev      *clipPath
}

type clip struct {
	extents RectangleInt
	path    *clipPath
	boxes   []box

	region      *region
	isRegion    bool
	embeddedBox [1]box
}
