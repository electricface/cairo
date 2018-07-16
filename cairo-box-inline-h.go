package cairo

func (box *box) set(p1, p2 *point) {
	box.p1 = *p1
	box.p2 = *p2
}

func boxFromIntegers(box *box, x, y, width, height int) {
	box.p1.x = fixedFromInt(x)
	box.p1.y = fixedFromInt(y)
	box.p2.x = fixedFromInt(x + width)
	box.p2.y = fixedFromInt(y + height)
}

func (box *box) addPoint(point *point) {
	if point.x < box.p1.x {
		box.p1.x = point.x
	} else if point.x > box.p2.x {
		box.p2.x = point.x
	}

	if point.y < box.p1.y {
		box.p1.y = point.y
	} else if point.y > box.p2.y {
		box.p2.y = point.y
	}
}

func (box *box) addBox(add *box) {
	if add.p1.x < box.p1.x {
		box.p1.x = add.p1.x
	}
	if add.p2.x > box.p2.x {
		box.p2.x = add.p2.x
	}

	if add.p1.y < box.p1.y {
		box.p1.y = add.p1.y
	}
	if add.p2.y > box.p2.y {
		box.p2.y = add.p2.y
	}
}

func (box *box) containsPoint(point *point) bool {
	return box.p1.x <= point.x && point.x <= box.p2.x &&
		box.p1.y <= point.y && point.y <= box.p2.y
}

func (box *box) isPixelAligned() bool {
	var f fixed
	f |= box.p1.x & fixedFracMask
	f |= box.p1.y & fixedFracMask
	f |= box.p2.x & fixedFracMask
	f |= box.p2.y & fixedFracMask
	return f == 0
}
