package rating

import "math"

//ComputeInitialVolatility is utils function for detamine volatility
//start and end is deviation. count is rating period count.
func ComputeInitialVolatility(start, end, count float64) float64 {
	return nthFloor(math.Sqrt((math.Pow(end/convartRate, 2)-math.Pow(start/convartRate, 2))/count), 6)
}
