package consensus

// CalcEMAUint64 takes a slice of uint64 values and returns the last EMA value,
// using a smoothing value encoded as a ratio of numerator and denominator and a
// specified SMA sampling window
//
// in energi PoS v2 consensus this is 1/3, generated by from
//
//     2/(params.SMAPeriod+1)
//
// This function assumes that there is no negative intervals between the
// samples, a monotonic timestamped blockchain being the use for this function.
// Use of data that is not known to be monotonic will have overflow values in
// any interval that computes as negative with signed integers.
func CalcEMAUint64(
	samples []uint64,
	numerator,
	denominator uint64,
	smaWindow uint64,
) (o uint64) {
	// nothing to do, nothing to do
	if len(samples) < 2 {
		return
	}
	// the result is derived from the intervals between so that it is 1 less
	// than the number input
	intervals := make([]uint64, len(samples)-1)
	// first we generate the raw intervals between the samples
	for i := 1; i < len(samples); i++ {
		intervals[i-1] = samples[i] - samples[i-1]
	}
	sma := make([]uint64, len(intervals)-int(smaWindow))
	// next generate the Simple Moving Average with smaWindow window
	for i := 0; i < len(sma); i++ {
		for j := int(smaWindow) - 1; j >= 0; j-- {
			sma[i] += intervals[j+i]
		}
		sma[i] /= smaWindow
	}
	// then compute the EMA with the given smoothing ratio
	//
	// EMA = (closing price − previous day’s EMA) × smoothing constant as a
	// decimal * previous day’s EMA
	//
	// The last clause of the formula is equivalent to multiplying by a
	// fraction, such as 2/(5+1) as used in this difficulty adjustment
	// algorithm
	o = sma[0]
	for i := range sma {
		if i > 0 {
			o = sma[i]*numerator/denominator +
				sma[i-1]*(denominator-numerator)/denominator
		}
	}
	return
}
