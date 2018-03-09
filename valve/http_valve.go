package valve

import (
	"net/http"

	"github.com/throttled/throttled"
)

type keyFunc func(r *http.Request) string

type HTTPValve struct {
	Valve
	key keyFunc
}

func NewHTTPValve(store throttled.Store, key keyFunc, schedule *Schedule) HTTPValve {
	return HTTPValve{NewValve(store, schedule), key}
}

func (r HTTPValve) RateLimit(req *http.Request) (bool, throttled.RateLimitResult, error) {
	return r.Valve.RateLimit(r.key(req), 1)
}
