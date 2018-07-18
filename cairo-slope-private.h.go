package cairo

func (s *slope) init(a, b *point) {
	s.dx = b.x - a.x
	s.dy = b.y - a.y
}

func (a *slope) equal(b *slope) bool {
	return int64(a.dy)*int64(b.dx) == int64(b.dy)*int64(a.dx)
}

func (a *slope) backwards(b *slope) bool {
	return (int64(a.dx)*int64(b.dx))+(int64(a.dy)*int64(b.dy)) < 0
}
