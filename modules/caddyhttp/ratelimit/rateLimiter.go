package ratelimit

import (
	"math"
	"time"
)

var rateLimiter *rateLimitOptions

// rateLimitOptions stores options detailing how rate limiting should be applied,
// as well as the current and previous window's hosts:requestCount mapping
type rateLimitOptions struct {
	// window length for request rate checking (>= 5 minutes)
	windowLength time.Duration

	// max request that should be processed per host in a given windowLength
	maxRequests int64

	// current window's request count per host
	currentWindow *requestCountTracker

	// previous window's request count per host
	previousWindow *requestCountTracker
}

// setupRateLimiter sets up the package-level variable `rateLimiter`,
// and starts the auto-window refresh process
func setupRateLimiter(windowLength time.Duration, maxRequests int64) (rl *rateLimitOptions) {
	rl = &rateLimitOptions{
		windowLength:   windowLength,
		maxRequests:    maxRequests,
		currentWindow:  newRequestCountTracker(windowLength),
		previousWindow: &requestCountTracker{},
	}

	go func() { // automatic shuffling of request count tracking windows
		for {
			time.Sleep(rl.currentWindow.endTime.Sub(time.Now()))
			rl.refreshWindows()
		}
	}()

	return
}

// refreshWindows() checks if currentWindow has reached its expiry time, and if it has,
// moves currentWindow to previousWindow, and re-initialises currentWindow
func (rl *rateLimitOptions) refreshWindows() (didRefresh bool) {
	if rl.currentWindow.endTime.Before(time.Now()) {
		rl.previousWindow = rl.currentWindow
		rl.currentWindow = newRequestCountTracker(rl.windowLength)

		didRefresh = true
	}

	return
}

// requestShouldBlock checks whether the request from a given host name should block,
// and increments the request counter for the hostName first
// will block if current request would push the hostName over the blocking threshold
func (rl *rateLimitOptions) requestShouldBlock(hostName string) (shouldBlock bool) {
	rl.currentWindow.addRequestForHost(hostName)                     // increment request counter for host
	return rl.getInterpolatedRequestCount(hostName) > rl.maxRequests // check if they now are above the request limit
}

// getInterpolatedRequestCount gets an interpolated request count for a specified hostName
// Always considers requests across the given windowLength
// More details: https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
//
// For example say given a case where:
// 	windowLength is 20 minutes
// 	current window started 10 minutes ago
// 	requestCount would be 0.5 * currentWindowRequests + 0.5 * previousWindowRequests
func (rl rateLimitOptions) getInterpolatedRequestCount(hostName string) (requestCount int64) {
	now := time.Now()

	// calculate fraction of request that went in the current and previous windows
	currentWindowFraction := now.Sub(rl.currentWindow.startTime).Seconds() / rl.windowLength.Seconds()
	previousWindowFraction := 1 - currentWindowFraction // thankfully this one's a bit easier to calculate!

	requestCount += int64(math.Round(
		float64(rl.currentWindow.getRequestCountForHost(hostName)) *
			currentWindowFraction))
	requestCount += int64(math.Round(
		float64(rl.previousWindow.getRequestCountForHost(hostName)) *
			previousWindowFraction))

	return
}
