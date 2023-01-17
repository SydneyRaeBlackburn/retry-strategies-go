package retry

import "time"

type Backoff interface {
	nextBackoff() time.Duration
	reset()
	Retry(func() interface{}) error
}
