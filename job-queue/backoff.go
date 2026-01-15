package main

import (
	"math"
	"time"
)

func calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(
		float64(BaseBackoff) * math.Pow(2, float64(attempt-1)),
	)

	if backoff > MaxBackoff {
		return MaxBackoff
	}
	return backoff
}
