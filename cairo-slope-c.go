package cairo

func (a *slope) compare(b *slope) int {
	adyBdx := int64(a.dy) * int64(b.dx)
	bdyAdx := int64(b.dy) * int64(a.dx)
	var cmp int

	cmp = int64Cmp(adyBdx, bdyAdx)
	if cmp != 0 {
		return cmp
	}

	/* special-case zero vectors.  the intended logic here is:
	 * zero vectors all compare equal, and more positive than any
	 * non-zero vector.
	 */

	if a.dx == 0 && a.dy == 0 && b.dx == 0 && b.dy == 0 {
		return 0
	}
	if a.dx == 0 && a.dy == 0 {
		return 1
	}
	if b.dx == 0 && b.dy == 0 {
		return -1
	}

	/* Finally, we're looking at two vectors that are either equal or
	 * that differ by exactly pi. We can identify the "differ by pi"
	 * case by looking for a change in sign in either dx or dy between
	 * a and b.
	 *
	 * And in these cases, we eliminate the ambiguity by reducing the angle
	 * of b by an infinitesimally small amount, (that is, 'a' will
	 * always be considered less than 'b').
	 */
	if (a.dx^b.dx) < 0 || (a.dy^b.dy) < 0 {
		if a.dx > 0 || (a.dx == 0 && a.dy > 0) {
			return -1
		}
		return 1
	}

	/* Finally, for identical slopes, we obviously return 0. */
	return 0
}
