package cairo

type Cairo struct {
	status Status
	Backend *Backend
}

func (cr *Cairo) Status() Status {
	return cr.status
}

type Backend struct {

}




