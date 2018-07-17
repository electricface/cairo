package cairo

import (
	"container/list"
	"log"
)

func (boxes *boxes) init() {
	boxes.status = StatusSuccess
	boxes.limits = nil
	boxes.numBoxes = 0

	boxes.chunks = list.New()
	boxes.chunks.PushBack(makeBoxesChunk(32))

	boxes.isPixelAligned = true
}

func (boxes *boxes) initFromRectangle(x, y, w, h int) {
	boxes.init()
	front := boxes.chunks.Front()
	firstChunk := front.Value.(boxesChunk)

	var box0 box
	boxFromIntegers(&box0, x, y, w, h)
	firstChunk.boxes = append(firstChunk.boxes, box0)
	boxes.numBoxes = 1
}

func (boxes *boxes) initWithClip(clip *clip) {
	boxes.init()
	if clip != nil {
		boxes.doLimit(clip.boxes)
	}
}

func (boxes *boxes) initForArray(array []box) {
	numBoxes := len(array)
	var n int
	boxes.status = StatusSuccess
	boxes.limits = nil
	boxes.numBoxes = numBoxes

	boxes.chunks.PushBack(array)

	for n = 0; n < numBoxes; n++ {
		if !array[n].p1.x.IsInteger() ||
			!array[n].p1.y.IsInteger() ||
			!array[n].p2.x.IsInteger() ||
			!array[n].p2.y.IsInteger() {
			break
		}
	}
	boxes.isPixelAligned = n == numBoxes
}

func (boxes *boxes) doLimit(limits []box) {
	var n int
	boxes.limits = limits
	numLimits := len(limits)
	if numLimits == 0 {
		return
	}
	boxes.limit = limits[0]
	for n = 1; n < numLimits; n++ {
		if limits[n].p1.x < boxes.limit.p1.x {
			boxes.limit.p1.x = limits[n].p1.x
		}

		if limits[n].p1.y < boxes.limit.p1.y {
			boxes.limit.p1.y = limits[n].p1.y
		}

		if limits[n].p2.x > boxes.limit.p2.x {
			boxes.limit.p2.x = limits[n].p2.x
		}

		if limits[n].p2.y > boxes.limit.p2.y {
			boxes.limit.p2.x = limits[n].p2.x
		}

		if limits[n].p2.y > boxes.limit.p2.y {
			boxes.limit.p2.y = limits[n].p2.y
		}
	}
}

func (boxes *boxes) addInternal(box0 *box) {
	if boxes.status != 0 {
		return
	}

	back := boxes.chunks.Back()
	chunk := back.Value.(boxesChunk)
	if len(chunk.boxes) == cap(chunk.boxes) {
		// chunk full
		newCap := cap(chunk.boxes) * 2
		chunk = makeBoxesChunk(newCap)
		boxes.chunks.PushBack(chunk)
	}
	chunk.boxes = append(chunk.boxes, *box0)
	boxes.numBoxes++

	if boxes.isPixelAligned {
		boxes.isPixelAligned = box0.isPixelAligned()
	}
}

func (boxes *boxes) add(antialias Antialias, box0 *box) Status {
	var b box

	if antialias == AntialiasNone {
		b.p1.x = box0.p1.x.roundDown()
		b.p1.y = box0.p1.y.roundDown()
		b.p2.x = box0.p2.x.roundDown()
		b.p2.y = box0.p2.y.roundDown()
		box0 = &b
	}

	if box0.p1.y == box0.p2.y {
		return StatusSuccess
	}

	if box0.p1.x == box0.p2.x {
		return StatusSuccess
	}

	if len(boxes.limits) != 0 {
		var p1, p2 point
		reversed := false

		/* support counter-clockwise winding for rectangular tessellation */
		if box0.p1.x < box0.p2.x {
			p1.x = box0.p1.x
			p2.x = box0.p2.x
		} else {
			p2.x = box0.p1.x
			p1.x = box0.p2.x
			reversed = !reversed
		}

		if p1.x >= boxes.limit.p2.x || p2.x <= boxes.limit.p1.x {
			return StatusSuccess
		}

		if box0.p1.y < box0.p2.y {
			p1.y = box0.p1.y
			p2.y = box0.p2.y
		} else {
			p2.y = box0.p1.y
			p1.y = box0.p2.y
			reversed = !reversed
		}

		if p1.y >= boxes.limit.p2.y || p2.y <= boxes.limit.p1.y {
			return StatusSuccess
		}

		for _, limits := range boxes.limits {
			var _box box
			var _p1, _p2 point

			if p1.x >= limits.p2.x || p2.x <= limits.p1.x {
				continue
			}
			if p1.y >= limits.p2.y || p2.y <= limits.p1.y {
				continue
			}

			/* Otherwise, clip the box to the limits. */
			_p1 = p1
			if _p1.x < limits.p1.x {
				_p1.x = limits.p1.x
			}
			if _p1.y < limits.p1.y {
				_p1.y = limits.p1.y
			}

			_p2 = p2
			if _p2.x > limits.p2.x {
				_p2.x = limits.p2.x
			}
			if _p2.y > limits.p2.y {
				_p2.y = limits.p2.y
			}

			if _p2.y <= _p1.y || _p2.x <= _p1.x {
				continue
			}

			_box.p1.y = _p1.y
			_box.p2.y = _p2.y
			if reversed {
				_box.p1.x = _p2.x
				_box.p2.x = _p1.x
			} else {
				_box.p1.x = _p1.x
				_box.p2.x = _p2.x
			}
			boxes.addInternal(&_box)
		}

	} else {
		boxes.addInternal(box0)
	}

	return boxes.status
}

func (boxes *boxes) extents(box0 *box) {
	var chunk boxesChunk
	var b box
	var i int

	if boxes.numBoxes == 0 {
		box0.p1.x = 0
		box0.p1.y = 0
		box0.p2.x = 0
		box0.p2.y = 0
		return
	}

	front := boxes.chunks.Front()
	b = front.Value.(boxesChunk).boxes[0]
	for elem := front; elem != nil; elem = elem.Next() {
		chunk = elem.Value.(boxesChunk)
		for _, bi := range chunk.boxes {
			if bi.p1.x < b.p1.x {
				b.p1.x = bi.p1.x
			}

			if bi.p1.y < b.p1.y {
				b.p1.y = bi.p1.y
			}

			if bi.p2.x > b.p2.x {
				b.p2.x = bi.p2.x
			}

			if bi.p2.y > b.p2.y {
				b.p2.y = bi.p2.y
			}
		}
	}
	*box0 = b
}

func (boxes *boxes) clear() {
	front := boxes.chunks.Front()
	chunk := front.Value.(boxesChunk)
	chunk.boxes = chunk.boxes[:0]
	boxes.chunks = list.New()
	boxes.chunks.PushBack(chunk) // re-use old chunk
	boxes.numBoxes = 0
	boxes.isPixelAligned = true
}

func (boxes *boxes) toArray() []box {
	ret := make([]box, boxes.numBoxes)

	j := 0
	for elem := boxes.chunks.Front(); elem != nil; elem = elem.Next() {
		chunk := elem.Value.(boxesChunk)
		j += copy(ret[j:], chunk.boxes)
	}
	return ret
}

func (boxes *boxes) fini() {

}

func (boxes *boxes) forEachBox(fn func(*box) bool) bool {
	for elem := boxes.chunks.Front(); elem != nil; elem = elem.Next() {
		chunk := elem.Value.(boxesChunk)
		for _, b := range chunk.boxes {
			if !fn(&b) {
				// fn return false then stop
				return false
			}
		}
	}
	return true
}

type boxRenderer struct {
	base  spanRenderer
	boxes *boxes
}

// 考虑让 boxRenderer 实现 spanRenderer 接口，rendRows 就用 spanToBoxes 实现。

func spanToBoxes(abstractRenderer interface{}, y, h int, spans []halfOpenSpan) Status {
	r := abstractRenderer.(boxRenderer)
	status := StatusSuccess
	var box0 box

	if len(spans) == 0 {
		return StatusSuccess
	}

	box0.p1.y = fixedFromInt(y)
	box0.p2.y = fixedFromInt(y + h)

	numSpans := len(spans)
	for {
		if spans[0].coverage != 0 {
			box0.p1.x = fixedFromInt(int(spans[0].x))
			box0.p2.x = fixedFromInt(int(spans[1].x))
			status = r.boxes.add(AntialiasDefault, &box0)
		}
		numSpans--
		if !(numSpans > 1 && status == StatusSuccess) {
			break
		}
		spans = spans[1:]
	}
	return status
}

func rasterisePolygonToBoxes(polygon *polygon, fillRule FillRule, boxes *boxes) Status {
	// TODO:
	var renderer boxRenderer
	var converter scanConverter
	var status Status
	var r RectangleInt

	boxRoundToRectangle(&polygon.extents, &r)
	converter = monoScanConverterCreate(r.X, r.Y,
		r.X+r.Width,
		r.Y+r.Height,
		fillRule)
	status = monoScanConverterAddPolygon(converter, polygon)
	if status != 0 {
		goto cleanupConverter
	}

	renderer.boxes = boxes
	// TODO:
	// renderer.base.renderRows = spanToBoxes

	status = converter.generate(converter, &renderer.base)

cleanupConverter:
	converter.destroy()
	return status
}

func (boxes *boxes) debugPrint() {
	var extents box

	boxes.extents(&extents)
	log.Printf("boxes x %d: (%f, %f) x (%f, %f)\n",
		boxes.numBoxes,
		extents.p1.x.toDouble(),
		extents.p1.y.toDouble(),
		extents.p2.x.toDouble(),
		extents.p2.y.toDouble(),
	)

	for elem := boxes.chunks.Front(); elem != nil; elem = elem.Next() {
		chunk := elem.Value.(boxesChunk)
		for i, b := range chunk.boxes {
			log.Printf("  box[%d]: (%f, %f), (%f, %f)\n", i,
				b.p1.x.toDouble(),
				b.p1.y.toDouble(),
				b.p2.x.toDouble(),
				b.p2.y.toDouble(),
			)
		}
	}
}
