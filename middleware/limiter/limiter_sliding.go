package limiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
)

type SlidingWindow struct{}

// New creates a new sliding window middleware handler
func (SlidingWindow) New(cfg Config) fiber.Handler {
	var (
		// Limiter variables
		mux        = &sync.RWMutex{}
		expiration = uint64(cfg.Expiration.Seconds())
	)

	// Create manager to simplify storage operations ( see manager.go )
	manager := newManager(cfg.Storage)

	// Update timestamp every second
	utils.StartTimeStampUpdater()

	// Return new handler
	return func(c fiber.Ctx) error {
		// Generate maxRequests from generator, if no generator was provided the default value returned is 5
		maxRequests := cfg.MaxFunc(c)

		// Don't execute middleware if Next returns true or if the max is 0
		if (cfg.Next != nil && cfg.Next(c)) || maxRequests == 0 {
			return c.Next()
		}

		// Get key from request
		key := cfg.KeyGenerator(c)

		// Lock entry
		mux.Lock()

		// Get entry from pool and release when finished
		e := manager.get(c, key)

		// Get timestamp
		ts := uint64(utils.Timestamp())

		// Set expiration if entry does not exist
		if e.exp == 0 {
			e.exp = ts + expiration
		} else if ts >= e.exp {
			// The entry has expired, handle the expiration.
			// Set the prevHits to the current hits and reset the hits to 0.
			e.prevHits = e.currHits

			// Reset the current hits to 0.
			e.currHits = 0

			// Check how much into the current window it currently is and sets the
			// expiry based on that, otherwise this would only reset on
			// the next request and not show the correct expiry.
			elapsed := ts - e.exp
			if elapsed >= expiration {
				e.exp = ts + expiration
			} else {
				e.exp = ts + expiration - elapsed
			}
		}

		// Increment hits
		e.currHits++

		// Calculate when it resets in seconds
		resetInSec := e.exp - ts

		// weight = time until current window reset / total window length
		weight := float64(resetInSec) / float64(expiration)

		// rate = request count in previous window - weight + request count in current window
		rate := int(float64(e.prevHits)*weight) + e.currHits

		// Calculate how many hits can be made based on the current rate
		remaining := cfg.Max - rate

		// Update storage. Garbage collect when the next window ends.
		// |--------------------------|--------------------------|
		//               ^            ^               ^          ^
		//              ts         e.exp   End sample window   End next window
		//               <------------>
		// 				   Reset In Sec
		// resetInSec = e.exp - ts - time until end of current window.
		// duration + expiration = end of next window.
		// Because we don't want to garbage collect in the middle of a window
		// we add the expiration to the duration.
		// Otherwise after the end of "sample window", attackers could launch
		// a new request with the full window length.
		manager.set(c, key, e, time.Duration(resetInSec+expiration)*time.Second) //nolint:gosec // Not a concern

		// Unlock entry
		mux.Unlock()

		// Check if hits exceed the cfg.Max
		if remaining < 0 {
			// Return response with Retry-After header
			// https://tools.ietf.org/html/rfc6584
			if !cfg.DisableHeaders {
				c.Set(fiber.HeaderRetryAfter, strconv.FormatUint(resetInSec, 10))
			}

			// Call LimitReached handler
			return cfg.LimitReached(c)
		}

		// Continue stack for reaching c.Response().StatusCode()
		// Store err for returning
		err := c.Next()

		// Get the effective status code from either the error or response
		statusCode := getEffectiveStatusCode(c, err)

		// Check for SkipFailedRequests and SkipSuccessfulRequests
		if (cfg.SkipSuccessfulRequests && statusCode < fiber.StatusBadRequest) ||
			(cfg.SkipFailedRequests && statusCode >= fiber.StatusBadRequest) {
			// Lock entry
			mux.Lock()
			e = manager.get(c, key)
			e.currHits--
			remaining++
			manager.set(c, key, e, cfg.Expiration)
			// Unlock entry
			mux.Unlock()
		}

		// We can continue, update RateLimit headers
		if !cfg.DisableHeaders {
			c.Set(xRateLimitLimit, strconv.Itoa(maxRequests))
			c.Set(xRateLimitRemaining, strconv.Itoa(remaining))
			c.Set(xRateLimitReset, strconv.FormatUint(resetInSec, 10))
		}

		return err
	}
}
