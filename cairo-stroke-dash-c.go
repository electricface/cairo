package cairo

func (dash *strokerDash) start() {
	var offset float64
	on := true
	var i int

	if !dash.dashed {
		return
	}

	offset = dash.dashOffset

	/* We stop searching for a starting point as soon as the
	offset reaches zero.  Otherwise when an initial dash
	segment shrinks to zero it will be skipped over. */
	for offset > 0.0 && offset >= dash.dashes[i] {
		offset -= dash.dashes[i]
		on = !on
		i++
		if i == len(dash.dashes) {
			i = 0
		}
	}

	dash.dashIndex = uint(i)
	dash.dashOn = on
	dash.dashStartsOn = on
	dash.dashRemain = dash.dashes[i] - offset
}

func (dash *strokerDash) step(step float64) {
	dash.dashRemain -= step
	if dash.dashRemain < fixedErrorDouble {
		dash.dashIndex++
		if int(dash.dashIndex) == len(dash.dashes) {
			dash.dashIndex = 0
		}
		dash.dashOn = !dash.dashOn
		dash.dashRemain += dash.dashes[dash.dashIndex]
	}
}

func (dash *strokerDash) init(style *strokeStyle) {
	dash.dashed = style.dash != nil
	if !dash.dashed {
		return
	}

	dash.dashes = style.dash
	dash.dashOffset = style.dashOffset
	dash.start()
}
