package valve_test

import (
	"testing"
	"time"

	"github.com/keratin/throttled-valve/valve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/throttled/throttled"
)

func TestRateLimit(t *testing.T) {
	store, err := NewTestStore()
	require.NoError(t, err)

	schedule := valve.NewSchedule(4,
		// burst: 4, period: 15 seconds, count: 4
		valve.Entry{Rate: throttled.PerMin(4)},
		// burst: 6, period: 30 seconds, count: 2
		valve.Entry{Rate: throttled.PerMin(2), Delay: 1 * time.Minute},
	)
	v := valve.NewValve(store, schedule)

	store.clock = time.Unix(0, 0)

	var limited bool
	var result throttled.RateLimitResult

	// use up the initial burst
	// t=0-14.99
	limited, result, err = v.RateLimit("foo", 5)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, time.Duration(-1), result.RetryAfter)

	// hit the first limit
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.True(t, limited)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 15*time.Second, result.RetryAfter)

	// t=15-29.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 5, result.Limit)

	// t=30-44.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 5, result.Limit)

	// t=45-59.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 7, result.Limit)

	// t=60-74.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 7, result.Limit)

	// t=75-89.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.True(t, limited)
	assert.Equal(t, 7, result.Limit)

	// t=90-104.99
	store.clock = store.clock.Add(15 * time.Second)
	limited, result, err = v.RateLimit("foo", 1)
	require.NoError(t, err)
	assert.False(t, limited)
	assert.Equal(t, 7, result.Limit)
}
