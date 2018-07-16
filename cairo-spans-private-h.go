package cairo

import (
	"container/list"
	"unsafe"
)

const (
	spansUnitCoverageBits = 8
	spansUnitCoverage     = (1 << spansUnitCoverageBits) - 1
)

type halfOpenSpan struct {
	x        int32
	coverage uint8
	inverse  uint8
}

type spanRenderer interface {
	getStatus() Status
	setError(status Status)

	destroy()
	renderRows(abstractRenderer interface{}, y, height int, coverages []halfOpenSpan)
	finish(abstractRenderer interface{}) Status
}

type scanConverter interface {
	destroy()
	generate(abstractRenderer interface{}, renderer spanRenderer) Status

	getStatus() Status
	setError(status Status)
}

const stackBufferSize = 512 * unsafe.Sizeof(int(0))

type rectangularScanConverter struct {
	status  Status
	extents box

	chunks        *list.List // elem type is rectangularScanConverterChunk
	buf           [stackBufferSize]byte
	numRectangles int
}

func (*rectangularScanConverter) destroy() {
	panic("implement me")
}

func (*rectangularScanConverter) generate(abstractRenderer interface{}, renderer spanRenderer) Status {
	panic("implement me")
}

func (*rectangularScanConverter) getStatus() Status {
	panic("implement me")
}

func (*rectangularScanConverter) setError(status Status) {
	panic("implement me")
}

type rectangularScanConverterChunk struct {
	data []byte
}

type botorScanConverter struct {
	status   Status
	extents  box
	fillRule FillRule

	xMin, xMax int

	chunks   *list.List // elem type is botorScanConverterChunk
	buf      [stackBufferSize]byte
	numEdges int
}

func (*botorScanConverter) destroy() {
	panic("implement me")
}

func (*botorScanConverter) generate(abstractRenderer interface{}, renderer spanRenderer) Status {
	panic("implement me")
}

func (*botorScanConverter) getStatus() Status {
	panic("implement me")
}

func (*botorScanConverter) setError(status Status) {
	panic("implement me")
}

type botorScanConverterChunk struct {
	data []byte
}
