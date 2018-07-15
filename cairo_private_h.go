package cairo

type Cairo struct {
	status  Status
	backend backend
}

func (cr *Cairo) Status() Status {
	return cr.status
}
