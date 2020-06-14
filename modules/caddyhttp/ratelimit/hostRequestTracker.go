package ratelimit

import "time"

type requestCountTracker struct {
	requestCount map[string]int64 // If 9,223,372,036,854,775,807 requests isn't enough...
	startTime    time.Time
	endTime      time.Time
}

// newRequestCountTracker returns a pointer to a blank initialised requestCountTracker
func newRequestCountTracker() *requestCountTracker {
	return &requestCountTracker{
		requestCount: map[string]int64{},
		startTime:    time.Now(),
		endTime:      time.Now().Add(rateLimiter.windowLength),
	}
}

// addRequestForHost adds to the request counter for specified host name
func (rct *requestCountTracker) addRequestForHost(hostName string) {
	rct.requestCount[hostName] += 1
}

// getRequestCounterForHost gets the request count for a given host name
func (rct requestCountTracker) getRequestCountForHost(hostName string) (requestCount int64) {
	return rct.requestCount[hostName]
}
