package cairo

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
