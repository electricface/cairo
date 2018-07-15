package cairo

func boxSet(box *box, p1, p2 *point) {
	box.p1 = *p1
	box.p2 = *p2
}

func boxFromIntegers(box *box, x, y, w, h int) {
	box.p1.x = fixedFromInt(x)
	box.p1.y = fixedFromInt(y)

}
