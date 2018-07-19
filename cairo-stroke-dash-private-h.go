package cairo

type strokerDash struct {
	dashed       bool
	dashIndex    uint
	dashOn       bool
	dashStartsOn bool
	dashRemain   float64
	dashOffset   float64
	dashes       []float64
}
