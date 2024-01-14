package utils

import "math"

// MaxInt returns the larger of x or y for type int.
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}

	return y
}

// MinInt returns the smaller of x or y for type int.
func MinInt(x int, y int) int {
	if x > y {
		return y
	}

	return x
}

// MaxFloat32 returns the larger of x or y for type float32.
func MaxFloat32(x float32, y float32) float32 {
	if x > y {
		return x
	}

	return y
}

// MinFloat32 returns the smaller of x or y for type float32.
func MinFloat32(x float32, y float32) float32 {
	if x > y {
		return y
	}

	return x
}

// DegToRad converts degrees to radians.
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// RadToDeg converts radians to degrees.
func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}
