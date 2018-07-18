package cairo

func uint64DivRem(num, den uint64) uQuoRem64 {
	var qr uQuoRem64
	qr.quo = num / den
	qr.rem = num % den
	return qr
}

func int64DivRem(num, den int64) quoRem64 {
	var uqr uQuoRem64
	var qr quoRem64

	if num < 0 {
		num = -num
	}
	if den < 0 {
		den = -den
	}

	uqr = uint64DivRem(uint64(num), uint64(den))

	if num < 0 {
		qr.rem = -int64(uqr.rem)
	} else {
		qr.rem = int64(uqr.rem)
	}

	if (num < 0) != (den < 0) {
		qr.quo = -int64(uqr.quo)
	} else {
		qr.quo = int64(uqr.quo)
	}
	return qr
}

func (v uint128) toUint64() uint64 {
	return v.lo
}

func (v uint128) toUint32() uint32 {
	return uint32(v.toUint64())
}

func (v uint128) isZero() bool {
	return v.hi == 0 && v.lo == 0
}

func (v int128) toInt64() int64 {
	return int64(v.lo)
}

func (v int128) toInt32() int32 {
	return int32(v.toInt64())
}

func uint64Cmp(a, b uint64) int {
	if a == b {
		return 0
	} else if a < b {
		return -1
	}
	return 1
}

func int64Cmp(a, b int64) int {
	if a == b {
		return 0
	} else if a < b {
		return -1
	}
	return 1
}
