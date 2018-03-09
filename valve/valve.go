package valve

import (
	"github.com/throttled/throttled"
)

// Valve is a RateLimiter implementation comprising multiple limiters that must each be queried and
// maintained on every unit of work. It is intended to be used in situations where sustained work is
// suspicious because it is sustained, so that the sustained rate can be throttled down to a crawl.
//
// Accomplishing this requires querying and maintaining each limiter on every call to RateLimit. A
// Valve with 3 limiters will take 3x longer to verify, which may be meaningful if checking a
// limiter requires network I/O to a centralized store.
type Valve struct {
	limiters []*throttled.GCRARateLimiter
}

// NewValve builds a Valve from a Schedule
func NewValve(store throttled.Store, schedule *Schedule) *Valve {
	var limiters []*throttled.GCRARateLimiter
	for _, r := range schedule.Rates {
		limiter, _ := throttled.NewGCRARateLimiter(store, *r)
		limiters = append(limiters, limiter)
	}
	return &Valve{limiters}
}

// RateLimit will return true (valve is limited) if any of its limiters are limited. The
// RateLimitResult will describe the limiter with the least remaining capacity.
func (q *Valve) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	var result throttled.RateLimitResult
	for idx, limiter := range q.limiters {
		limited, _result, err := limiter.RateLimit(key+string(idx), quantity)
		if err != nil || limited {
			return limited, _result, err
		}
		if idx == 0 || _result.Remaining <= result.Remaining {
			result = _result
		}
	}
	return false, result, nil
}
