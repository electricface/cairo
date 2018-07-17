package cairo

type fixed16_16 int32

type fixed32_32 int64

type fixed48_16 int64

type fixed64_64 int128
type fixed96_32 int128

const fixedBits = 32
const fixedFracBits = 8

type fixed int32

type fixedUnsigned uint32

type point struct {
	x, y fixed
}
