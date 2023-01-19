# Retry Strategies

This package allows an api call to be retried based on the chosen backoff strategy. 

## Backoff Strategies

### Constant
`wait_duration = constant + jitter`

- TODO

### Linear
`wait_duration = last_duration + delta + jitter`

- TODO

### Fibonacci
- TODO

### Quadratic
- TODO

### Exponential
`wait_duration = intial_delay_interval * scaling_factor^(attempt) + jitter`

```go
// can be reused
eb := b.NewExponentialBackoff(b.ExponentialBackoff{
    //InitialInterval: 250 * time.Millisecond,
    //MaxInterval: 10 * time.Second,
    //MaxAttempts: 10, 
    //ScalingFactor: 1,
})

// declare outside Retry() for scope
var tf *tfexec.Terraform
var stdErr error

// top level error will return a failure on retry i.e maxAttempts/maxInterval reached
err := eb.Retry(func() interface{} {
    // stdErr returns errors from the external api call 
    // this error is checked in Retry()
    // if nil, will break retry loop; else, log error and continue
    tf, stdErr = s.init(ctx.GetContext(), logger, req.GetResourceId(), tfFolder)
    return stdErr
})
if err != nil {
    return err
}

// do something with tf
```

### Polynomial
- TODO

## Jitter
- TODO

## TODO
- Implement other strategies (linear, fibonacci, quadratic, polunomial, etc.)
- Unit tests
- Better Documentation
    - Use cases, example