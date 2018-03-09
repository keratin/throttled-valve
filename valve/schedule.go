package valve

import (
	"time"

	"github.com/throttled/throttled"
)

// Schedule is configuration struct for NewValve. It comprises multiple RateQuotas that have been
// calculated by NewSchedule to work together in a constricting fashion.
type Schedule struct {
	Rates []*throttled.RateQuota
}

// Entry is a configuration struct for NewSchedule. It describes a Rate that will only take effect
// after some Delay in maximally sustained traffic according to the previous rate.
type Entry struct {
	Rate  throttled.Rate
	Delay time.Duration
}

// NewSchedule constructs a Schedule from the set of Entries provided. The Entries are converted to
// throttled.RateQuota structs with burst values calculated to allow the previous rate in the
// Schedule to be maximally sustained during the delay.
func NewSchedule(burst int, entries ...Entry) *Schedule {
	previousBurst := burst
	var previousPeriod time.Duration

	rates := []*throttled.RateQuota{}
	for _, entry := range entries {
		thisPeriod := periodOfRate(entry.Rate)

		var consumedDuringDelay int
		if int(previousPeriod) == 0 {
			consumedDuringDelay = previousBurst
		} else {
			consumedDuringDelay = previousBurst + int(entry.Delay)/int(previousPeriod)
		}
		refilledDuringDelay := int(entry.Delay) / int(thisPeriod)

		thisRate := &throttled.RateQuota{
			MaxBurst: consumedDuringDelay - refilledDuringDelay,
			MaxRate:  entry.Rate,
		}

		rates = append(rates, thisRate)
		previousBurst = thisRate.MaxBurst
		previousPeriod = thisPeriod
	}

	return &Schedule{rates}
}

// throttled/throttled does not make period public, but it does provide enough to calculate it out
func periodOfRate(r throttled.Rate) time.Duration {
	count, interval := r.Quota()
	return interval / time.Duration(count)
}
