package cairo

import (
	"fmt"
	"io"
)

var clipAll clip

func (c *clip) pathCreate() *clipPath {
	cp := new(clipPath)

	cp.prev = c.path
	c.path = cp
	return cp
}

func (cp *clipPath) destroy() {
}

func clipCreate() *clip {
	var c *clip
	c = new(clip)
	// TODO: clip pool
	c.extents = unboundedRectangle
	c.path = nil
	c.boxes = nil
	c.region = nil
	c.isRegion = false
	return c
}

func (c *clip) destroy() {
	if c == nil || c.isAllClipped() {
		return
	}

	if c.path != nil {
		c.path.destroy()
	}

	c.region.destroy()
	// TODO: put c to clip pool
}

func (c *clip) copy() *clip {
	var cp *clip
	if c == nil || c.isAllClipped() {
		return c
	}
	cp = clipCreate()

	if c.path != nil {
		cp.path = c.path
	}

	if len(c.boxes) != 0 {
		if len(c.boxes) == 1 {
			cp.boxes = cp.embeddedBox[:]
		} else {
			cp.boxes = make([]box, len(c.boxes))
			if cp.boxes == nil {
				return cp.setAllClipped()
			}
		}
		copy(cp.boxes, c.boxes)
	}

	cp.extents = c.extents
	cp.region = c.region
	cp.isRegion = c.isRegion

	return cp
}

func (c *clip) copyPath() *clip {
	if c == nil || c.isAllClipped() {
		return c
	}

	if len(c.boxes) == 0 {
		panic("assert failed len(c.boxes) != 0")
	}

	cp := clipCreate()
	cp.extents = c.extents
	if c.path != nil {
		cp.path = c.path
	}
	return cp
}

func (c *clip) copyRegion() *clip {
	if c == nil || c.isAllClipped() {
		return c
	}

	if len(c.boxes) == 0 {
		panic("assert failed len(c.boxes) != 0")
	}

	cp := clipCreate()
	cp.extents = c.extents

	if len(c.boxes) == 1 {
		cp.boxes = cp.embeddedBox[:]
	} else {
		cp.boxes = make([]box, len(c.boxes))
		if cp.boxes == nil {
			return cp.setAllClipped()
		}
	}

	for i := 0; i < len(c.boxes); i++ {
		cp.boxes[i].p1.x = c.boxes[i].p1.x.floor()
		cp.boxes[i].p1.y = c.boxes[i].p1.y.floor()
		cp.boxes[i].p2.x = c.boxes[i].p2.x.ceil()
		cp.boxes[i].p2.y = c.boxes[i].p2.y.ceil()
	}

	cp.region = c.region
	cp.isRegion = true
	return cp
}

func (c *clip) intersectPath(path *pathFixed, fillRule FillRule, tolerance float64,
	antialias Antialias) *clip {

	if c.isAllClipped() {
		return c
	}

	/* catch the empty clip path */
	if path.fillIsEmpty() {
		return c.setAllClipped()
	}

	var box box
	if path.isBox(&box) {
		if antialias == AntialiasNone {
			box.p1.x = box.p1.x.roundDown()
			box.p1.y = box.p1.y.roundDown()
			box.p2.x = box.p2.x.roundDown()
			box.p2.y = box.p2.x.roundDown()
		}
		return c.intersectBox(&box)
	}

	if path.fillIsRectilinear() {
		return c.intersectRectilinearPath(path, fillRule, antialias)
	}

	var extents RectangleInt
	path.approximateClipExtents(&extents)
	if extents.Width == 0 || extents.Height == 0 {
		return c.setAllClipped()
	}

	c = c.intersectRectangle(&extents)
	if c.isAllClipped() {
		return c
	}

	clipPath := c.pathCreate()
	if clipPath == nil {
		return c.setAllClipped()
	}

	status := clipPath.path.initCopy(path)
	if status != 0 {
		return c.setAllClipped()
	}

	clipPath.fillRule = fillRule
	clipPath.tolerance = tolerance
	clipPath.antialias = antialias

	if c.region != nil {
		c.region.destroy()
		c.region = nil
	}
	c.isRegion = false
	return c
}

func (c *clip) intersectClipPath(clipPath *clipPath) *clip {
	if clipPath.prev != nil {
		c = c.intersectClipPath(clipPath.prev)
	}

	return c.intersectPath(&clipPath.path, clipPath.fillRule, clipPath.tolerance,
		clipPath.antialias)
}

func (c *clip) intersectClip(other *clip) *clip {
	if c.isAllClipped() {
		return c
	}

	if other == nil {
		return c
	}

	if c == nil {
		return other.copy()
	}

	if other.isAllClipped() {
		return c.setAllClipped()
	}

	if !c.extents.intersect(&other.extents) {
		return c.setAllClipped()
	}

	if len(other.boxes) != 0 {
		var boxes boxes
		boxes.initForArray(other.boxes)
		c = c.intersectBoxes(&boxes)
	}

	if !c.isAllClipped() {
		if other.path != nil {
			if c.path == nil {
				c.path = other.path
			} else {
				c = c.intersectClipPath(other.path)
			}
		}
	}

	if c.region != nil {
		c.region.destroy()
		c.region = nil
	}
	c.isRegion = false

	return c
}

func (clipA *clip) equal(clipB *clip) bool {
	/* are both all-clipped or no-clip? */
	if clipA == clipB {
		return true
	}

	/* or just one of them? */
	if clipA == nil || clipB == nil ||
		clipA.isAllClipped() ||
		clipB.isAllClipped() {
		return false
	}

	/* We have a pair of normal clips, check their contents */
	if len(clipA.boxes) != len(clipB.boxes) {
		return false
	}

	for i := 0; i < len(clipA.boxes); i++ {
		if clipA.boxes[i] != clipB.boxes[i] {
			return false
		}
	}

	cpA := clipA.path
	cpB := clipB.path
	for cpA != nil && cpB != nil {
		if cpA == cpB {
			return true
		}

		/* XXX compare reduced polygons? */

		if cpA.antialias != cpB.antialias {
			return false
		}

		if cpA.tolerance != cpB.tolerance {
			return false
		}

		if cpA.fillRule != cpB.fillRule {
			return false
		}

		if !cpA.path.equal(cpB.path) {
			return false
		}

		cpA = cpA.prev
		cpB = cpB.prev
	}

	return cpA == nil && cpB == nil
}

func (c *clip) pathCopyWithTranslation(otherPath *clipPath, fx, fy int) *clip {
	if otherPath.prev != nil {
		c = c.pathCopyWithTranslation(otherPath.prev, fx, fy)
	}

	if c.isAllClipped() {
		return c
	}

	clipPath := c.pathCreate()
	if clipPath == nil {
		return c.setAllClipped()
	}
	status := clipPath.path.initCopy(&otherPath.path)
	if status != 0 {
		return c.setAllClipped()
	}

	clipPath.path.translate(fx, fy)

	clipPath.fillRule = otherPath.fillRule
	clipPath.tolerance = otherPath.tolerance
	clipPath.antialias = otherPath.antialias

	return c
}

func (c *clip) translate(tx, ty int) *clip {
	if c == nil || c.isAllClipped() {
		return c
	}

	if tx == 0 && ty == 0 {
		return c
	}

	fx := fixedFromInt(tx)
	fy := fixedFromInt(ty)

	for i := 0; i < len(c.boxes); i++ {
		c.boxes[i].p1.x += fx
		c.boxes[i].p2.x += fx
		c.boxes[i].p1.y += fy
		c.boxes[i].p2.y += fy
	}

	c.extents.X += tx
	c.extents.Y += ty

	if c.path == nil {
		return c
	}

	clipPath := c.path
	c.path = nil
	c = c.pathCopyWithTranslation(clipPath, int(fx), int(fy))
	clipPath.destroy()

	return c
}

func (path *pathFixed) addBox(box *box) Status {
	status := path.moveTo(box.p1.x, box.p1.y)
	if status != 0 {
		return status
	}

	status = path.lineTo(box.p2.x, box.p1.y)
	if status != 0 {
		return status
	}

	status = path.lineTo(box.p2.x, box.p2.y)
	if status != 0 {
		return status
	}

	status = path.lineTo(box.p1.x, box.p2.y)
	if status != 0 {
		return status
	}

	return path.closePath()
}

func (path *pathFixed) initFromBoxes(boxes *boxes) Status {
	var status Status

	path.init()
	if boxes.numBoxes == 0 {
		return StatusSuccess
	}

	for elem := boxes.chunks.Front(); elem != nil; elem = elem.Next() {
		chunk := elem.Value.(boxesChunk)
		for _, box := range chunk.boxes {
			status = path.addBox(&box)
			if status != 0 {
				path.fini()
				return status
			}
		}
	}
	return StatusSuccess
}

func (c *clip) transform(m *Matrix) *clip {
	var cp *clip
	if c == nil || c.isAllClipped() {
		return c
	}

	if m.isTranslation() {
		return c.translate(int(m.X0), int(m.Y0))
	}

	cp = clipCreate()

	if len(c.boxes) != 0 {
		var path pathFixed
		var boxes boxes
		boxes.initForArray(c.boxes)
		path.initFromBoxes(&boxes)
		path.transform(m)

		cp = cp.intersectPath(&path, FillRuleWinding, 0.1, AntialiasDefault)
		path.fini()
	}

	if c.path != nil {
		cp = cp.intersectClipPathTransformed(c.path, m)
	}

	c.destroy()
	return cp
}

func (c *clip) copyWithTranslation(tx, ty int) *clip {
	var cp *clip

	if c == nil || c.isAllClipped() {
		return c
	}

	if tx == 0 && ty == 0 {
		return c.copy()
	}
	cp = clipCreate()
	if cp == nil {
		return cp.setAllClipped()
	}

	fx := fixedFromInt(tx)
	fy := fixedFromInt(ty)

	if len(c.boxes) != 0 {
		if len(c.boxes) == 1 {
			cp.boxes = cp.embeddedBox[:]
		} else {
			cp.boxes = make([]box, len(c.boxes))
			if cp.boxes == nil {
				return cp.setAllClipped()
			}
		}

		for i := 0; i < len(c.boxes); i++ {
			cp.boxes[i].p1.x = c.boxes[i].p1.x + fx
			cp.boxes[i].p2.x = c.boxes[i].p2.x + fx
			cp.boxes[i].p1.y = c.boxes[i].p1.y + fy
			cp.boxes[i].p2.y = c.boxes[i].p2.y + fy
		}
	}

	cp.extents = c.extents
	cp.extents.X += tx
	cp.extents.Y += ty

	if c.path == nil {
		return cp
	}

	return cp.pathCopyWithTranslation(c.path, int(fx), int(fy))
}

func (c *clip) containsExtents(extents *compoisteRectangles) bool {
	var rect *RectangleInt
	if extents.isBounded {
		rect = &extents.bounded
	} else {
		rect = &extents.unbounded
	}
	return c.containsRectangle(rect)
}

func (c *clip) debugPrint(w *io.Writer) {
	if c == nil {
		fmt.Fprintln(w, "no clip")
		return
	}

	if c.isAllClipped() {
		fmt.Fprintln(w, "clip: all-clipped")
		return
	}

	fmt.Fprintln(w, "clip:")
	fmt.Fprintf(w, "  extents: (%d, %d) x (%d, %d), is-region? %v",
		c.extents.X, c.extents.Y,
		c.extents.Width, c.extents.Height,
		c.isRegion)

	fmt.Fprintf(w, "  num_boxes = %d\n", len(c.boxes))
	for i, b := range c.boxes {
		fmt.Fprintf(w, "  [%d] = (%f, %f), (%f, %f)\n", i,
			b.p1.x.toDouble(),
			b.p1.y.toDouble(),
			b.p2.x.toDouble(),
			b.p2.y.toDouble())
	}

	if c.path != nil {
		clipPath := c.path
		for {
			fmt.Fprintf(w, "path: aa=%d, tolerance=%f, rule=%d: ",
				clipPath.antialias, clipPath.tolerance, clipPath.fillRule)
			clipPath.path.debugPrint(w)
			fmt.Fprintln(w)

			clipPath = clipPath.prev
			if clipPath == nil {
				break
			}
		}
	}
}

func (c *clip) getExtents() *RectangleInt {
	if c == nil {
		return &unboundedRectangle
	}
	if c.isAllClipped() {
		return &emptyRectangle
	}
	return &c.extents
}

var rectanglesNil = RectangleList{
	Status: StatusNoMemory,
}

var rectanglesNotRepresentable = RectangleList{
	Status: StatusClipNotRepresentable,
}

func clipIntRectToUser(gstate *gstate, clipRect *RectangleInt, userRect *Rectangle) bool {
	var isTight bool

	x1 := float64(clipRect.X)
	y1 := float64(clipRect.Y)
	x2 := float64(clipRect.X + clipRect.Width)
	y2 := float64(clipRect.Y + clipRect.Height)

	gstateBackendToUserRectangle(gstate, &x1, &y1, &x2, &y2, &isTight)

	userRect.X = x1
	userRect.Y = y1
	userRect.Width = x2 - x1
	userRect.Height = y2 - y1
	return isTight
}

func rectangleListCreateInError(status Status) *RectangleList {
	var list *RectangleList
	if status == StatusNoMemory {
		return &rectanglesNil
	}
	if status == StatusClipNotRepresentable {
		return &rectanglesNotRepresentable
	}

	list = new(RectangleList)
	if list == nil {
		status = StatusNoMemory.error()
		return &rectanglesNil
	}
	list.Status = status
	list.Rectangles = nil
	return list
}

func (c *clip) copyRectangleList(gstate *gstate) *RectangleList {
	var list *RectangleList
	var rectangles []Rectangle
	var region *region
	var nRects int

	if c == nil {
		return rectangleListCreateInError(StatusClipNotRepresentable.error())
	}

	if c.isAllClipped() {
		goto DONE
	}

	if !c.getIsRegion() {
		return rectangleListCreateInError(StatusClipNotRepresentable.error())
	}

	region = c.getRegion()
	if region == nil {
		return rectangleListCreateInError(StatusNoMemory.error())
	}

	nRects = region.numRectangles()
	if nRects != 0 {
		rectangles = make([]Rectangle, nRects)
		if rectangles == nil {
			return rectangleListCreateInError(StatusNoMemory.error())
		}

		for i := 0; i < nRects; i++ {
			var clipRect RectangleInt
			region.getRectangle(i, &clipRect)
			if !clipIntRectToUser(gstate, &clipRect, &rectangles[i]) {
				rectangles = nil
				return rectangleListCreateInError(StatusClipNotRepresentable.error())
			}
		}
	}

DONE:
	list = new(RectangleList)
	if list == nil {
		rectangles = nil
		return rectangleListCreateInError(StatusNoMemory.error())
	}

	list.Status = StatusSuccess
	list.Rectangles = rectangles
	return list
}

func (list *RectangleList) destroy() {
	if list == nil || list == &rectanglesNil ||
		list == &rectanglesNotRepresentable {
		return
	}
	list.Rectangles = nil
}

func clipResetStaticData() {
}
