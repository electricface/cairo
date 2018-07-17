package cairo

type traps struct {
	status Status
	bounds box
	limits []box

	maybeRegion      bool
	hasIntersections bool
	isRectilinear    bool
	isRectangular    bool

	numTraps      int
	trapsSize     int
	traps         []trapzoid
	trapsEmbedded [16]trapzoid
}
