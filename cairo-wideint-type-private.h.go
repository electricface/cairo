package cairo

type uQuoRem64 struct {
	quo, rem uint64
}

type quoRem64 struct {
	quo, rem int64
}

type uint128 struct {
	lo, hi uint64
}

type int128 uint128

type uQuoRem128 struct {
	quo, rem uint128
}

type quoRem128 struct {
	quo, rem int128
}
