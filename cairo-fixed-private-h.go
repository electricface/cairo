package cairo

import "math"

func fixedFromInt(i int) fixed {
	return fixed(i << fixedFracBits)
}

const magicNumberFixed16_16 = 103079215104.0

const magicNumberFixed = (1 << (52 - fixedFracBits)) * 1.5

func fixedFromDouble(d float64) fixed {
	// TODO: need test
	d = d + magicNumberFixed
	return fixed(math.Float64bits(d) >> 32)
}

func fixedFrom26_6(i uint32) fixed {
	if fixedFracBits > 6 {
		return fixed(i << (fixedFracBits - 6))
	}
	return fixed(i >> uint(6-fixedFracBits))
}

func fixedFrom16_16(i uint32) fixed {
	if fixedFracBits > 16 {
		return fixed(1 << uint(fixedFracBits-16))
	}
	return fixed(1 >> (16 - fixedFracBits))
}

const (
	fixedOne         = fixed(1 << fixedFracBits)
	fixedOneDouble   = float64(1 << fixedFracBits)
	fixedEpsilon     = fixed(1)
	fixedErrorDouble = 1.0 / (2 * fixedOneDouble)

	fixedFracMask  = fixed((^fixedUnsigned(0)) >> (fixedFracBits - fixedFracBits))
	fixedWholeMask = ^fixedFracMask
)

func (f fixed) toDouble() float64 {
	return float64(f) / fixedOneDouble
}

func (f fixed) IsInteger() bool {
	return f&fixedFracMask == 0
}

func (f fixed) floor() fixed {
	return f &^ fixedFracMask
}

func (f fixed) ceil() fixed {
	return (f + fixedFracMask).floor()
}

func (f fixed) round() fixed {
	return (f + (fixedFracMask+1)/2).floor()
}

func (f fixed) roundDown() fixed {
	return (f + fixedFracMask/2).floor()
}

func (f fixed) integerPart() int {
	return int(f >> fixedFracBits)
}

func (f fixed) integerRound() int {
	return (f + (fixedFracMask+1)/2).integerPart()
}

func (f fixed) integerRoundDown() int {
	return (f + fixedFracMask/2).integerPart()
}

func (f fixed) fractionalPart() int {
	return int(f & fixedFracMask)
}

func (f fixed) integerFloor() int {
	if f >= 0 {
		return int(f >> fixedFracBits)
	}
	return int(-((-f - 1) >> fixedFracBits) - 1)
}

func (f fixed) integerCeil() int {
	if f > 0 {
		return int(((f - 1) >> fixedFracBits) + 1)
	}
	return int(-(-f >> fixedFracBits))
}

func (f fixed) to16_16() fixed16_16 {
	if fixedFracBits == 16 && fixedBits == 32 {
		return fixed16_16(f)
	} else if fixedFracBits > 16 {
		return fixed16_16(f >> uint(fixedFracBits-16))
	}
	var x fixed16_16
	if (f >> fixedFracBits) < math.MinInt16 {
		x = math.MinInt32
	} else if (f >> fixedFracBits) > math.MaxInt16 {
		x = math.MaxInt32
	} else {
		x = fixed16_16(f << (16 - fixedFracBits))
	}

	return x
}

func fixed16_16FromDouble(d float64) fixed16_16 {
	d = d + magicNumberFixed16_16
	return fixed16_16(math.Float64bits(d) >> 32)
}

func (f fixed16_16) floor() int {
	if f >= 0 {
		return int(f >> 16)
	}
	return int(-((-f - 1) >> 16) - 1)
}

func (f fixed16_16) toDouble() float64 {
	return float64(f) / float64(1<<16)
}

func (a fixed) mul(b fixed) fixed {
	var temp int64 = int32x32_64mul(a, b)
	return int64ToInt32(int64Rsl(temp, fixedFracBits))
}

func (a fixed) mulDiv(b, c fixed) fixed {
	var ab int64 = int32x32_64mul(a, b)
	var c64 int64 = int32ToInt64(c)
	var int64ToInt32(int64Divrem(ab, c64).quo)
}

func (a fixed) mulDivFloor(b, c fixed) fixed {
	return int64_32div(int32x32_64mul(a, b), c)
}

func edgeComputeIntersectionYForX(p1, p2 *point, x fixed) fixed {
	var y , dx fixed

	if x == p1.x {
		return p1.y
	}
	if x == p2.x {
		return p2.y
	}

	y = p1.y
	dx = p2.x - p1.x
	if dx != 0 {
		y += (x - p1.x).mulDivFloor(p2.y - p1.y, dx)
	}
	return y
}

func edgeComputeIntersectionXForY(p1,p2 *point, y fixed) fixed {
	var x, dy fixed
	if y == p1.y {
		return p1.x
	}
	if y == p2.y {
		return p2.x
	}

	x = p1.x
	dy = p2.y - p1.y
	if dy != 0 {
		x += (y - p1.y).mulDivFloor(p2.x - p1.x, dy)
	}
	return x
}

func segmentIntersection(seg1p1, seg1p2, seg2p1, seg2p2, intersection *point) bool {
	var denominator, uA, uB float64
	var seg1dx, seg1dy, seg2dx, seg2dy, segStartDx, segStartDy float64

	seg1dx = (seg1p2.x - seg1p1.x).toDouble()
	seg1dy = (seg1p2.y - seg1p1.y).toDouble()
	seg2dx = (seg2p2.x - seg2p1.x).toDouble()
	seg2dy = (seg2p2.y - seg2p1.y).toDouble()
	denominator = (seg2dy * seg1dx) - (seg2dx * seg1dy)
	if denominator == 0 {
		return false
	}

	segStartDx = (seg1p1.x - seg2p1.x).toDouble()
	segStartDy = (seg1p1.y - seg2p1.y).toDouble()
	uA = ((seg2dx * segStartDy) - (seg2dy * segStartDx)) / denominator
	uB = ((seg1dx * segStartDy) - (seg1dy * segStartDx)) / denominator

	if uA <= 0 || uA >= 1 || uB <= 0 || uB >= 1 {
		return false
	}

	intersection.x = seg1p1.x + fixedFromDouble(uA * seg1dx)
	intersection.y = seg1p1.y + fixedFromDouble(uA * seg1dy)
	return true
}
