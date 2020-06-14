package ratelimit

import "time"

type requestCountTracker struct {
	requestCount map[string]int64 // If 9,223,372,036,854,775,807 requests isn't enough...
	startTime    time.Time
	endTime      time.Time
}

func newRequestCountTracker() *requestCountTracker {
	return &requestCountTracker{
		requestCount: map[string]int64{},
		startTime:    time.Now(),
		endTime:      time.Now().Add(rateLimiter.windowLength),
	}
}
