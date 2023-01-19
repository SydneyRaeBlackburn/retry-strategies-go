package retry

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type LinearBackoff struct {
	MaxInterval    time.Duration // max time to wait between a retry attempt. default is 60s
	currentAttempt int
	MaxAttempts    int // max retry attempts before failing. default is 10
	Delta          int // constant increase (must be greater than 0). default is 1
}

const (
	linearMaxIntervalDefault = 60 * time.Second
	linearMaxAttemptsDefault = 3
	deltaDefault             = 2
)

var (
	linearMaxInterval time.Duration
	linearMaxAttempts int
	delta             int
)

func init() {
	linearMaxInterval = linearMaxIntervalDefault
	linearMaxAttempts = linearMaxAttemptsDefault
	delta = deltaDefault
}

// NewLinearBackoff returns a Backoff interface of type LinearBackoff
func NewLinearBackoff(lb LinearBackoff) Backoff {
	// set defaults
	if lb.MaxInterval.Seconds() <= 0 {
		lb.MaxInterval = linearMaxIntervalDefault
	}
	linearMaxInterval = lb.MaxInterval

	if lb.MaxAttempts <= 0 {
		lb.MaxAttempts = linearMaxAttemptsDefault
	}
	linearMaxAttempts = lb.MaxAttempts

	if lb.Delta <= 1 {
		lb.Delta = deltaDefault
	}
	delta = lb.Delta

	return &lb
}

// NextBackoff calculates the next wait interval for Retry()
func (b *LinearBackoff) nextBackoff() time.Duration {
	// jitter adds some amount of randomness to the backoff to spread the retries around in time
	// this prevents successive collision
	jitter := time.Duration(1000-rand.Intn(2000)) * time.Millisecond
	return time.Duration((b.currentAttempt-1)*b.Delta+b.Delta)*time.Second + jitter
}

// Reset will set the backoff values back to the defaults
func (b *LinearBackoff) reset() {
	b.MaxInterval = linearMaxInterval
	b.currentAttempt = 0
	b.MaxAttempts = linearMaxAttempts
	b.Delta = delta
}

//	 Retry will continuously execute a given function over an increasing amount of time
//		until a max interval or max amount of attempts has been reached,
//		or a successful call has been achieved
func (b *LinearBackoff) Retry(f func() interface{}) error {
	for {
		// increment and check if maxAttempts have been reached
		b.currentAttempt++
		if b.currentAttempt > b.MaxAttempts {
			err := errors.New("max retry attempts reached")
			b.reset() // reset backoff for next use
			return err
		}

		// get next wait time
		waitInterval := b.nextBackoff()

		// retry after current wait interval
		select {
		case <-time.After(waitInterval):
			// retry function call after wait time has elapsed
			err := f()
			if err != nil {
				fmt.Println(err, "..retrying after next interval")
				continue
			}

			// successful function call
			b.reset() // reset backoff for next use
			return nil
		case <-time.After(b.MaxInterval):
			err := errors.New("max retry interval reached")
			b.reset() // reset backoff for next use
			return err
		}
	}
}
