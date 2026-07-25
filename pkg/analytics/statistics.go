package analytics

import (
	"math"
)

func ComputeConfidenceScore(samples int, variance float64) float64 {
	if samples <= 0 {
		return 0.50
	}
	base := 0.70 + math.Min(0.25, float64(samples)/600.0)
	if variance > 0 {
		base -= math.Min(0.10, variance/100.0)
	}
	if base < 0.50 {
		base = 0.50
	}
	if base > 0.99 {
		base = 0.99
	}
	return base
}
