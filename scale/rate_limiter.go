package main

import (
	"sync"
	"time"
)

type rateLimiter struct {
	burstRate     int
	rate          int
	currentTokens int
	mu            sync.Mutex
	ticker        *time.Ticker
	quitChan      <-chan struct{}
	isClosed      bool
}

func (r *rateLimiter) New(rate, burstRate int) *rateLimiter {
	rl := &rateLimiter{
		burstRate:     burstRate,
		rate:          rate,
		currentTokens: rate,
		quitChan:      make(<-chan struct{}),
		ticker:        time.NewTicker(time.Second / time.Duration(rate)),
		isClosed:      false,
	}

	go rl.refill()
	return rl
}

func (r *rateLimiter) refill() {
	for {
		select {
		case <-r.ticker.C:
			r.mu.Lock()
			if r.burstRate > r.currentTokens {
				r.currentTokens++
			}
			r.mu.Unlock()
		case <-r.quitChan:
			r.ticker.Stop()
			return
		}
	}
}

func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentTokens > 0 {
		r.currentTokens--
		return true

	}
	return false
}

func (r *rateLimiter) Close() {

}
func rateLimitDriver() {

}
