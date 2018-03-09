package valve_test

import (
	"testing"
	"time"

	"github.com/keratin/throttled/valve"
	"github.com/stretchr/testify/assert"
	"github.com/throttled/throttled"
)

func TestNewSchedule(t *testing.T) {
	schedule := valve.NewSchedule(4,
		valve.Entry{Rate: throttled.PerMin(10)},
		valve.Entry{Rate: throttled.PerMin(3), Delay: 2 * time.Minute},
		valve.Entry{Rate: throttled.PerMin(1), Delay: 10 * time.Minute},
	)

	expected := &valve.Schedule{[]*throttled.RateQuota{
		{MaxRate: throttled.PerMin(10), MaxBurst: 4},
		{MaxRate: throttled.PerMin(3), MaxBurst: 18},
		{MaxRate: throttled.PerMin(1), MaxBurst: 38},
	}}

	assert.Equal(t, expected, schedule)
}
