package retry

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type ExponentialBackoff struct {
	InitialInterval time.Duration // initial wait time. default is .5s
	MaxInterval     time.Duration // max time to wait between a retry attempt. default is 60s
	currentAttempt  int
	MaxAttempts     int // max retry attempts before failing. default is 10
	ScalingFactor   int // base multiplier of the power(must be greater than 1). default is 2
}

const (
	initialIntervalDefault = 500 * time.Millisecond
	maxIntervalDefault     = 60 * time.Second
	maxAttemptsDefault     = 3
	scalingFactorDefault   = 2
)

var (
	initialInterval time.Duration
	maxInterval     time.Duration
	maxAttempts     int
	scalingFactor   int
)

func init() {
	initialInterval = initialIntervalDefault
	maxInterval = maxIntervalDefault
	maxAttempts = maxAttemptsDefault
	scalingFactor = scalingFactorDefault
}

// NewExponentialBackoff returns a Backoff interface of type ExponentialBackoff
func NewExponentialBackoff(eb ExponentialBackoff) Backoff {
	// set defaults
	if eb.InitialInterval.Seconds() <= 0 {
		eb.InitialInterval = initialIntervalDefault
	}
	initialInterval = eb.InitialInterval

	if eb.MaxInterval.Seconds() <= 0 {
		eb.MaxInterval = maxIntervalDefault
	}
	maxInterval = eb.MaxInterval

	if eb.MaxAttempts <= 0 {
		eb.MaxAttempts = maxAttemptsDefault
	}
	maxAttempts = eb.MaxAttempts

	if eb.ScalingFactor <= 1 {
		eb.ScalingFactor = scalingFactorDefault
	}
	scalingFactor = eb.ScalingFactor

	return &eb
}

// NextBackoff calculates the next wait interval for Retry()
func (b *ExponentialBackoff) nextBackoff() time.Duration {
	// jitter adds some amount of randomness to the backoff to spread the retries around in time
	// this prevents successive collision
	jitter := time.Duration(1000-rand.Intn(2000)) * time.Millisecond
	return time.Duration(
		b.InitialInterval.Seconds()*math.Pow(
			float64(b.ScalingFactor),
			float64(b.currentAttempt),
		))*time.Second + jitter
}

// Reset will set the backoff values back to the defaults
func (b *ExponentialBackoff) reset() {
	b.InitialInterval = initialInterval
	b.MaxInterval = maxInterval
	b.currentAttempt = 0
	b.MaxAttempts = maxAttempts
	b.ScalingFactor = scalingFactor
}

// Retry will continuously execute a given function over an increasing amount of time
//	until a max interval or max amount of attempts has been reached,
//	or a successful call has been achieved
// 	ex. with a maxInterval of 60s, if the next backoff interval is > 60s then an err will return
func (b *ExponentialBackoff) Retry(f func() interface{}) error {
	for {
		// increment and check if maxAttempts have been reached
		b.currentAttempt++
		if b.currentAttempt > b.MaxAttempts {
			err := errors.New("max retry attempts reached")
			// logger.Error(err)
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
				// logger.WithError(err).Error("error from external api call.. retrying after next interval")
				fmt.Println(err, "..retrying after next interval")
				continue
			}

			// successful function call
			b.reset() // reset backoff for next use
			return nil
		case <-time.After(b.MaxInterval):
			err := errors.New("max retry interval reached")
			// logger.Error(err)
			b.reset() // reset backoff for next use
			return err
		}
	}
}
