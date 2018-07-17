package cairo

import "os"

type edgeT struct {
	next, prev, right *edgeT
	x, top            fixed
	dir               int
}

type rectangleT struct {
	left, right edgeT
	top, bottom int32
}

func pqParentIndex(i int) int {
	return i >> 1
}

const pqFirstEntry = 1

func pqLeftChildIndex(i int) int {
	return i << 1
}

type sweepLine struct {
	rectangles     []*rectangleT
	stop           []*rectangleT
	head, tail     edgeT
	insert, cursor *edgeT
	currentY       int32
	lastY          int32
	insertX        int32
	fillRule       FillRule
	doTraps        bool
	container      interface{}
	unwind         jmpBuf
}

type jmpBuf struct {
}

func (traps *traps) dump() {
	if os.Getenv("CAIRO_DEBUG_TRAPS") == "" {
		return
	}

}

/*
dump_traps (cairo_traps_t *traps, const char *filename)
{
    FILE *file;
    int n;

    if (getenv ("CAIRO_DEBUG_TRAPS") == NULL)
	return;

    file = fopen (filename, "a");
    if (file != NULL) {
	for (n = 0; n < traps->num_traps; n++) {
	    fprintf (file, "%d %d L:(%d, %d), (%d, %d) R:(%d, %d), (%d, %d)\n",
		     traps->traps[n].top,
		     traps->traps[n].bottom,
		     traps->traps[n].left.p1.x,
		     traps->traps[n].left.p1.y,
		     traps->traps[n].left.p2.x,
		     traps->traps[n].left.p2.y,
		     traps->traps[n].right.p1.x,
		     traps->traps[n].right.p1.y,
		     traps->traps[n].right.p2.x,
		     traps->traps[n].right.p2.y);
	}
	fprintf (file, "\n");
	fclose (file);
    }
}
*/
