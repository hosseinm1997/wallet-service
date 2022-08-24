package services

import (
	"arvan-wallet-service/types/exceptions"
	"errors"
	"github.com/sony/gobreaker"
	"sync"
	"sync/atomic"
	"time"
)

type state struct {
	counter uint32
	lastErr error
}

type mapState map[string]state

type creditCodeStateHolder struct {
	cb    *gobreaker.CircuitBreaker
	mutex sync.Mutex
	m     atomic.Value
}

var (
	shOnce sync.Once
	sh     *creditCodeStateHolder
)

func getStateHolder() *creditCodeStateHolder {
	//mux.Lock()
	//defer mux.Unlock()
	if sh == nil {
		shOnce.Do(func() {
			cfg := gobreaker.Settings{
				// When to flush counters' int the Closed state
				Interval: 5 * time.Second,
				// Time to switch from Open to Half-open
				Timeout: 2 * time.Second,
				// Function with check when to switch from Closed to Open
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
					return counts.ConsecutiveFailures >= 5 && failureRatio >= 0.7
				},
				IsSuccessful: func(err error) bool {
					if err == nil {
						return true
					}
					if errors.Is(err, exceptions.CreditCodeLimitationError{}) {
						return true
					}

					return false
				},
			}

			sh = &creditCodeStateHolder{}
			sh.cb = gobreaker.NewCircuitBreaker(cfg)
			sh.m.Store(make(mapState))
		})
	}

	return sh
}

func (c *creditCodeStateHolder) readValue(key string) *state {
	ms := c.m.Load().(mapState)

	if val, ok := ms[key]; ok {
		return &val
	}

	return nil
}

func (c *creditCodeStateHolder) storeValue(key string, val uint32, err error) {
	m1 := c.m.Load().(mapState) // load current value of the data structure
	m2 := make(mapState)        // create a new value
	for k, v := range m1 {
		m2[k] = v // copy all data from the current object to the new one
	}
	m2[key] = state{counter: val, lastErr: err} // do the update that we need
	c.m.Store(m2)                               // atomically replace the current object with the new one
	// At this point all new readers start working with the new version.
	// The old version will be garbage collected once the existing readers
	// (if any) are done with it.
}

func (c *creditCodeStateHolder) addValue(key string, val uint32, err error) {
	c.mutex.Lock() // synchronize with other potential writers
	defer c.mutex.Unlock()
	m1 := c.m.Load().(mapState) // load current value of the data structure
	m2 := make(mapState)        // create a new value
	for k, v := range m1 {
		m2[k] = v // copy all data from the current object to the new one
	}
	m2[key] = state{counter: m2[key].counter + val, lastErr: err} // do the update that we need
	c.m.Store(m2)                                                 // atomically replace the current object with the new one
	// At this point all new readers start working with the new version.
	// The old version will be garbage collected once the existing readers
	// (if any) are done with it.
}
