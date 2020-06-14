package ratelimit

import (
	"time"
)

var rateLimiter rateLimitOptions

// rateLimitOptions stores options detailing how rate limiting should be applied
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
func (rl *rateLimitOptions) setupRateLimiter(windowLength time.Duration, maxRequests int64) {
	rl.windowLength = windowLength
	rl.maxRequests = maxRequests

	go func() { // automatic shuffling of request count tracking windows
		for {
			time.Sleep(time.Now().Sub(rl.currentWindow.endTime))
			rl.refreshWindows()
		}
	}()
}

// refreshWindows() checks if currentWindow has reached its expiry time, and if it has,
// moves currentWindow to previousWindow, and re-initialises currentWindow
func (rl *rateLimitOptions) refreshWindows() (didRefresh bool) {
	if rl.currentWindow.endTime.Before(time.Now()) {

		rl.currentWindow, rl.previousWindow = newRequestCountTracker(), rl.currentWindow
		didRefresh = true
	}

	return
}
