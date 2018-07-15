package cairo

type backendType int

const (
	backendTypeDefault backendType = iota
	backendTypeSkia
)

type backend interface {
	getType() backendType
	destroy(cr *Cairo)
	getOriginTarget(cr *Cairo) *Surface
	getCurrentTarget(cr *Cairo) *Surface
	save(cr *Cairo) Status
	restore(cr *Cairo) Status
	pushGroup(cr *Cairo, content Content) Status
	popGroup(cr *Cairo) *Pattern
	setSourceRgba(cr *Cairo, red, green, blue, alpha float64) Status
	setSourceSurface(cr *Cairo, surface *Surface, x, y float64) Status
	setSource(cr *Cairo, source *Pattern)
	getSource(cr *Cairo) *Pattern
	setAntialias(cr *Cairo, antialias Antialias) Status
	setDash(cr *Cairo, dashes []float64, offset float64) Status
	setFillRule(cr *Cairo, rule FillRule) Status
	setLineCap(cr *Cairo, lineCap LineCap) Status
	setLineJoin(cr *Cairo, lineJoin LineJoin) Status
	setLineWidth(cr *Cairo, lineWidth float64) Status
	setMiterLimit(cr *Cairo, limit float64) Status
	setOpacity(cr *Cairo, opacity float64) Status
	setOperator(cr *Cairo, op Operator) Status
	setTolerance(cr *Cairo, tolerance float64) Status
	getAntialias(cr *Cairo) Antialias
	getDash(cr *Cairo) (dashes []float64, offset float64)
	getFillRule(cr *Cairo) FillRule
	getLineCap(cr *Cairo) LineCap
	getLineJoin(cr *Cairo) LineJoin
	getLineWidth(cr *Cairo) float64
	getMiterLimit(cr *Cairo) float64
	getOpacity(cr *Cairo) float64
	getOperator(cr *Cairo) Operator
	getTolerance(cr *Cairo) float64

	translate(cr *Cairo, tx, ty float64) Status
	scale(cr *Cairo, sx, sy float64) Status
	rotate(cr *Cairo, theta float64) Status
	transform(cr *Cairo, matrix *Matrix) Status
	setMatrix(cr *Cairo, matrix *Matrix) Status
	setIdentityMatrix(cr *Cairo) Status
	getMatrix(cr *Cairo) *Matrix

	userToDevice(cr *Cairo) (x, y float64)
	userToDeviceDistance(cr *Cairo) (x, y float64)
	deviceToUser(cr *Cairo) (x, y float64)
	deviceToUserDistance(cr *Cairo) (x, y float64)

	userToBackend(cr *Cairo) (x, y float64)
	userToBackendDistance(cr *Cairo) (x, y float64)
	backendToUser(cr *Cairo) (x, y float64)
	backendToUserDistance(cr *Cairo) (x, y float64)

	newPath(cr *Cairo) Status
	newSubPath(cr *Cairo) Status
	moveTo(cr *Cairo, x, y float64) Status
	relMoveTo(cr *Cairo, dx, dy float64) Status
	lineTo(cr *Cairo, x, y float64) Status
	relLineTo(cr *Cairo, dx, dy float64) Status
	curveTo(cr *Cairo, x1, y1, x2, y2, x3, y3 float64) Status
	relCurveTo(cr *Cairo, dx1, dy1, dx2, dy2, dx3, dy3 float64) Status
	arcTo(cr *Cairo, x1, y1, x2, y2, radius float64) Status
	relArcTo(cr *Cairo, dx1, dy1, dx2, dy2, radius float64) Status
	closePath(cr *Cairo) Status

	arc(cr *Cairo, xc, yc, radius, angle1, angle2 float64, forward bool) Status
	rectangle(cr *Cairo, x, y, width, height float64) Status

	pathExtents(cr *Cairo) (x1, y1, x2, y2 float64)
	hasCurrentPoint(cr *Cairo) bool
	getCurrentPoint(cr *Cairo) (has bool, x, y float64)

	copyPath(cr *Cairo) *Path
	copyPathFlat(cr *Cairo) *Path
	appendPath(cr *Cairo, path *Path) Status

	strokeToPath(cr *Cairo) Status

	clip(cr *Cairo) Status
	clipPreserve(cr *Cairo) Status
	inClip(cr *Cairo, x, y float64) (inside bool, status Status)
	clipExtents(cr *Cairo) (x1, y1, x2, y2 float64, status Status)
	resetClip(cr *Cairo) Status
	clipCopyRectangleList(cr *Cairo) RectangleList
	paint(cr *Cairo) Status
	paintWithAlpha(cr *Cairo, opacity float64) Status
	mask(cr *Cairo, pattern *Pattern) Status
	stroke(cr *Cairo) Status
	strokePreserve(cr *Cairo) Status
	inStroke(cr *Cairo) Status
	strokeExtents(cr *Cairo) (x1, y1, x2, y2 float64, status Status)

	fill(cr *Cairo) Status
	fillPreserve(cr *Cairo) Status
	inFill(cr *Cairo, x, y float64) (inside bool, status Status)
	fillExtents(cr *Cairo) (x1, y1, x2, y2 float64, status Status)

	setFontFace(cr *Cairo, fontFace *FontFace) Status
	getFontFace(cr *Cairo) *FontFace
	setFontSize(cr *Cairo, size float64) Status
	setFontMatrix(cr *Cairo, matrix *Matrix) Status
	getFontMatrix(cr *Cairo) *Matrix
	setFontOptions(cr *Cairo, options *fontOptions) Status
	getFontOptions(cr *Cairo) *fontOptions
	setScaledFont(cr *Cairo, scaledFont *scaledFont) Status
	getScaledFont(cr *Cairo) *scaledFont
	fontExtents(cr *Cairo) (extents *FontExtents, status Status)

	glyphs(cr *Cairo, glyphs []Glyph, info *glyphTextInfo)
	glyphPath(cr *Cairo, glyphs []Glyph) Status
	glyphExtents(cr *Cairo, glyphs []Glyph) (extents *TextExtents, status Status)

	copyPage(cr *Cairo) Status
	showPage(cr *Cairo) Status
}

type FontFace struct {
}

type scaledFont struct {
}

func (cr *Cairo) backendToUser() (x, y float64) {
	return cr.backend.backendToUser(cr)
}

func (cr *Cairo) backendToUserDistance() (x, y float64) {
	return cr.backend.backendToUserDistance(cr)
}

func (cr *Cairo) userToBackend() (x, y float64) {
	return cr.backend.userToBackend(cr)
}

func (cr *Cairo) userToBackendDistance() (x, y float64) {
	return cr.backend.userToBackendDistance(cr)
}
