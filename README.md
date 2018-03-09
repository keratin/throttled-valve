# Throttled-Valve

Extends [throttled/throttled](https://github.com/throttled/throttled) with a new RateLimiter that
can decrease the allowed rate of clients sustaining too much traffic for too long.

## Installation

```sh
go get -u github.com/keratin/throttled-valve
```

## Example

```go
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
http.ListenAndServe(":8080", loginLimiter.RateLimit(loginHandler))
```

## License

Copyright (c) 2018 Lance Ivy & contributors

Released under the MIT license.
