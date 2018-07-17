package cairo

func combsortNewgap(gap uint) uint {
	gap = 10 * gap / 13
	if gap == 9 || gap == 10 {
		gap = 11
	}
	if gap < 1 {
		gap = 1
	}
	return gap
}
