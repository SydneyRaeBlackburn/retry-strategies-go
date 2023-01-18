package retry

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type ConstantBackoff struct {
	Constant            time.Duration // max time to wait between a retry attempt. default is 60s
	currentAttempt      int
	ConstantMaxAttempts int // max retry attempts before failing. default is 10
}

const (
	constantDefault            = 5 * time.Second
	constantMaxAttemptsDefault = 10
)

var (
	constant            time.Duration
	constantMaxAttempts int
)

func init() {
	constant = constantDefault
	constantMaxAttempts = constantMaxAttemptsDefault
}

// NewConstantBackoff returns a Backoff interface of type ConstantBackoff
func NewConstantBackoff(cb ConstantBackoff) Backoff {
	// set defaults
	if cb.Constant.Seconds() <= 0 {
		cb.Constant = constantDefault
	}
	constant = cb.Constant

	if cb.ConstantMaxAttempts <= 0 {
		cb.ConstantMaxAttempts = constantMaxAttemptsDefault
	}
	constantMaxAttempts = cb.ConstantMaxAttempts

	return &cb
}

// NextBackoff calculates the next wait interval for Retry()
func (c *ConstantBackoff) nextBackoff() time.Duration {
	// jitter adds some amount of randomness to the backoff to spread the retries around in time
	// this prevents successive collision
	jitter := time.Duration(1000-rand.Intn(2000)) * time.Millisecond
	return c.Constant*time.Second + jitter
}

// Reset will set the backoff values back to the defaults
func (c *ConstantBackoff) reset() {
	c.Constant = constant
	c.currentAttempt = 0
	c.ConstantMaxAttempts = constantMaxAttempts
}

// Retry will continuously execute a given function over a
// constant amount of time until the max amount of attempts
// has been reached or a successful call has been achieved
func (c *ConstantBackoff) Retry(f func() interface{}) error {
	for {
		// increment and check if maxAttempts have been reached
		c.currentAttempt++
		if c.currentAttempt > c.ConstantMaxAttempts {
			err := errors.New("max retry attempts reached")
			c.reset() // reset backoff for next use
			return err
		}

		// get next wait time
		waitInterval := c.nextBackoff()

		// sleep for duration
		time.Sleep(waitInterval)

		// retry function call after wait time has elapsed
		err := f()
		if err != nil {
			fmt.Println(err, "..retrying after next interval")
			continue
		}

		// successful function call
		c.reset() // reset backoff for next use
		return nil
	}
}
