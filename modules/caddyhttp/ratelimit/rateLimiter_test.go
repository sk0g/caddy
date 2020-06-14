package ratelimit

import (
	"testing"
	"time"
)

func Test_rateLimitOptions_refreshWindows(t *testing.T) {
	t.Run("Should refresh", func(t *testing.T) {
		rateLimiter := rateLimitOptions{
			windowLength: 15 * time.Minute,
			currentWindow: &requestCountTracker{
				requestCount: map[string]int64{},
				startTime:    time.Now().Add(-35 * time.Minute),
				endTime:      time.Now().Add(-20 * time.Minute),
			},
		}

		if didRefresh := rateLimiter.refreshWindows(); !didRefresh {
			t.Errorf("Should have refreshed, but did not: %+v", rateLimiter)
		}
	})

	t.Run("Should not refresh", func(t *testing.T) {
		rateLimiter := rateLimitOptions{
			windowLength: 1 * time.Minute,
			currentWindow: &requestCountTracker{
				requestCount: map[string]int64{},
				startTime:    time.Now().Add(-30 * time.Minute),
				endTime:      time.Now().Add(30 * time.Minute),
			},
		}

		if didRefresh := rateLimiter.refreshWindows(); didRefresh {
			t.Errorf("Should not have refreshed, but did: %+v", rateLimiter)
		}
	})
}
