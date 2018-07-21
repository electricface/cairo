package cairo

func (path *pathFixed) isBox(box *box) bool {
	var buf *pathBuf
	if !path.fillIsRectilinear0 {
		return false
	}

	if !path.isQuad() {
		return false
	}

	buf = path.head()
	if pointsFromRect(buf.points) {
		canonicalBox(box, &buf.points[0], &buf.points[2])
		return true
	}

	return false
}
