package valve_test

import (
	"time"

	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

// from: https://github.com/throttled/throttled/blob/master/rate_test.go
type TestStore struct {
	store throttled.GCRAStore

	clock       time.Time
	failUpdates bool
}

func NewTestStore() (*TestStore, error) {
	mem, err := memstore.New(0)
	if err != nil {
		return nil, err
	}
	return &TestStore{store: mem}, err
}

func (ts *TestStore) GetWithTime(key string) (int64, time.Time, error) {
	v, _, e := ts.store.GetWithTime(key)
	return v, ts.clock, e
}

func (ts *TestStore) SetIfNotExistsWithTTL(key string, value int64, ttl time.Duration) (bool, error) {
	if ts.failUpdates {
		return false, nil
	}
	return ts.store.SetIfNotExistsWithTTL(key, value, ttl)
}

func (ts *TestStore) CompareAndSwapWithTTL(key string, old, new int64, ttl time.Duration) (bool, error) {
	if ts.failUpdates {
		return false, nil
	}
	return ts.store.CompareAndSwapWithTTL(key, old, new, ttl)
}
