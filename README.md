# Backoff

This package allows a service to retry an external api call based on the chosen backoff strategy. 

## Backoff Strategies

### Exponential with Jitter
`wait_duration = intial_delay_interval * scaling_factor^(attempt) + jitter`

```go
// example usage
import b "backoff"

// can be reused
eb := b.NewExponentialBackoff(b.ExponentialBackoff{
    //InitialInterval: 250 * time.Millisecond,
    //MaxInterval: 10 * time.Second,
    //MaxAttempts: 10, 
    //ScalingFactor: 1,
})

var tf *tfexec.Terraform
var stdErr error

// top level error will return a failure on retry i.e maxAttempts/maxInterval reached
err := eb.Retry(func() interface{} {
    // stdErr returns errors from external api call 
    // this error is checked in Retry()
    // if nil, will break retry loop; else, log error and continue
    tf, stdErr = s.init(ctx.GetContext(), logger, req.GetResourceId(), tfFolder)
    return stdErr
})
if err != nil {
    return logger.WrapError(err)
}

// do something with tf
```

## TODO
- implement other strategies (constant, linear, fibonacci, quadratic, polunomial, etc.)