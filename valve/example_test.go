package valve_test

import (
	"log"
	"net/http"
	"time"

	"github.com/keratin/throttled-valve/valve"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

// NOTE: keep in sync with README.md
func Example() {
	store, err := memstore.New(65536)
	if err != nil {
		log.Fatal(err)
	}

	// Valve will allow:
	// * 10/minute (1/6 seconds) with an additional burst of 4
	// * After two minutes: 2/minute (1/30 seconds)
	valve := valve.NewValve(store, valve.NewSchedule(4,
		valve.Entry{Rate: throttled.PerMin(10)},
		valve.Entry{Rate: throttled.PerMin(2), Delay: 2 * time.Minute},
	))

	// valve is compatible with throttled.HTTPRateLimiter
	loginLimiter := throttled.HTTPRateLimiter{
		RateLimiter: valve,
		VaryBy:      &throttled.VaryBy{RemoteAddr: true},
	}

	// you may wish to use a router and only apply the limit on specific endpoints
	var loginHandler http.HandlerFunc
	http.ListenAndServe(":8080", loginLimiter.RateLimit(loginHandler))
}
