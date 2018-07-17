package cairo

func (a *line) equal(b *line) bool {
	return a.p1.x == b.p1.x && a.p1.y == b.p1.y &&
		a.p2.x == b.p2.x && a.p2.y == b.p2.y
}
