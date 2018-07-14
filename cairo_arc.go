package cairo

import "math"

const maxFullCircles = 65535

func arcErrorNormalized(angle float64) float64 {
	return 2.0 / 27.0 * math.Pow(math.Sin(angle/ 4), 6) /
		math.Pow(math.Cos(angle/4),2)
}

func arcMaxAngleForToleranceNormalized(tolerance float64) float64 {
	var angle, error float64
	var i int

	/* Use table lookup to reduce search time in most cases. */
	table := []struct{
		angle, error float64
	}{
		{ math.Pi / 1.0,   0.0185185185185185036127 },
		{ math.Pi / 2.0,   0.000272567143730179811158 },
		{ math.Pi / 3.0,   2.38647043651461047433e-05 },
		{ math.Pi / 4.0,   4.2455377443222443279e-06 },
		{ math.Pi / 5.0,   1.11281001494389081528e-06 },
		{ math.Pi / 6.0,   3.72662000942734705475e-07 },
		{ math.Pi / 7.0,   1.47783685574284411325e-07 },
		{ math.Pi / 8.0,   6.63240432022601149057e-08 },
		{ math.Pi / 9.0,   3.2715520137536980553e-08 },
		{ math.Pi / 10.0,  1.73863223499021216974e-08 },
		{ math.Pi / 11.0,  9.81410988043554039085e-09 },
	}
	tableSize := len(table)
	for i = 0; i< tableSize; i++ {
		if table[i].error < tolerance {
			return table[i].angle
		}
	}

	i++
	for {
		angle = math.Pi / float64(i)
		i++
		error = arcErrorNormalized(angle)

		if error <= tolerance {
			break
		}
	}
	return angle
}

func arcSegmentsNeeded(angle, radius float64, ctm *Matrix, tolerance float64) int {
	var majorAxis, maxAngle float64
	majorAxis = ctm.TransformedCircleMajorAxis(radius)
	maxAngle = arcMaxAngleForToleranceNormalized (tolerance / majorAxis)

	return int(math.Ceil (math.Abs(angle) / maxAngle))
}

func (cr *Cairo) arcSegment(xc, yc, radius, angleA, angleB float64) {
	var rSinA,rCosA float64
	var rSinB, rCosB float64
	var h float64
	rSinA = radius * math.Sin(angleA)
	rCosA = radius * math.Cos(angleA)
	rSinB = radius * math.Sin(angleB)
	rCosB = radius * math.Cos(angleB)

	h = 4.0 / 3.0 * math.Tan( (angleB - angleA) / 4.0)

	cr.CurveTo(xc + rCosA - h * rSinA,
		yc + rSinA + h * rCosA,
			xc + rCosB + h * rSinB,
				yc + rSinB - h * rCosB,
					xc + rCosB,
						yc + rSinB)
}


func (cr *Cairo) arcInDirection(xc, yc , radius, angleMin, angleMax float64, dir direction) {
	if cr.Status() != 0 {
		return
	}

	if !(angleMax >= angleMin) {
		panic("assert failed")
	}

	if angleMax - angleMin > 2 * math.Pi * maxFullCircles {
		angleMax = math.Mod(angleMax - angleMin , 2 *math.Pi)
		angleMin = math.Mod(angleMin, 2*math.Pi)
		angleMax += angleMin + 2 * math.Pi * maxFullCircles
	}

	/* Recurse if drawing arc larger than pi */
	if angleMax - angleMin > math.Pi {
		angleMid := angleMin + (angleMax - angleMin) / 2.0
		if dir == directionForward {
			cr.arcInDirection(xc, yc, radius, angleMin, angleMid, dir)
			cr.arcInDirection(xc, yc, radius, angleMid, angleMax, dir)
		} else {
			cr.arcInDirection(xc, yc, radius, angleMid, angleMax, dir)
			cr.arcInDirection(xc, yc, radius, angleMin, angleMid, dir)
		}

	} else if angleMax != angleMin {
		var ctm Matrix
		var i, segments int
		var step float64

		cr.GetMatrix(&ctm)
		segments = arcSegmentsNeeded(angleMax - angleMin, radius, &ctm,
			cr.GetTolerance())
		step = (angleMax - angleMin) / float64(segments)
		segments = -1

		if dir == directionReverse {
			angleMin, angleMax = angleMax, angleMin
			step = - step
		}

		cr.LineTo(xc + radius * math.Cos(angleMin),
			yc + radius * math.Sin(angleMin))

		for i = 0; i < segments; i++{
			cr.arcSegment(xc, yc, radius, angleMin, angleMin + step)
			angleMin += step
		}

		cr.arcSegment(xc, yc, radius, angleMin, angleMax)
	} else {
		cr.LineTo(xc + radius * math.Cos(angleMin),
			yc + radius * math.Sin(angleMin))
	}
}

func (cr *Cairo) arcPath(xc, yc, radius, angle1, angle2 float64) {
	cr.arcInDirection(xc, yc, radius, angle1, angle2, directionForward)
}

func (cr *Cairo) arcPathNegative(xc, yc, radius, angle1, angle2 float64) {
	cr.arcInDirection(xc, yc, radius, angle1, angle2, directionReverse)
}
