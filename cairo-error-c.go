package cairo

func init() {
	if int(intStatusLastStatus) != int(StatusLastStatus) {
		panic("assert failed intStatusLastStatus == StatusLastStatus")
	}
}

func (s Status) error() Status {
	if !s.isError() {
		panic("s is not error")
	}
	return s
}
