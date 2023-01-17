package main

import (
	"errors"
	"fmt"
	"retry"
	"strconv"
	"time"
)

var num = 0

func main() {
	fmt.Println("HAPPY PATH WITH TIMES")
	// DECLARE NEW BACKOFF
	eb := retry.NewExponentialBackoff(retry.ExponentialBackoff{
		//InitialInterval: 250 * time.Millisecond,
		//MaxInterval: 10 * time.Second,
		MaxAttempts: 10, // comment out to fail on maxAttempts
		//ScalingFactor: 1,
	})

	start := time.Now()

	// EXAMPLE USAGE
	var v string     // declare return values outside Retry() for access
	var stdErr error // don't need to set this if a function only returns an error
	err := eb.Retry(func() interface{} {
		fmt.Println(time.Since(start)) // for testing
		v, stdErr = count(num)         // function to retry. needs to have an error returned
		return stdErr                  // need to return error from function here. this determines a retry. if nil, will not retry bc call was successful
	})
	if err != nil {
		fmt.Println(err)               // will be an error from the retry i.e maxAttempts/maxInterval reached
		fmt.Println(time.Since(start)) // for testing
	}

	// DO SOMETHING WITH RETURN VALUES IF YOU HAVE THEM
	fmt.Println(v)

	// reset
	num = 0
	fmt.Println("---------------------")
	fmt.Println("REUSE EXAMPLE")

	// CAN USE OLD BACKOFF WITH SAME DEFAULTS
	err = eb.Retry(func() interface{} {
		stdErr = countNoReturn(num)
		return stdErr
	})
	if err != nil {
		fmt.Println(err)
	}

	// reset
	num = 0
	fmt.Println("---------------------")
	fmt.Println("MAX ATTEMPTS FAILURE EXAMPLE")

	// OR LESS CODE WITH NEW INIT
	err = retry.NewExponentialBackoff(retry.ExponentialBackoff{}).Retry(func() interface{} {
		stdErr = countNoReturn(num)
		return stdErr
	})
	if err != nil {
		fmt.Println(err)
	}

	// reset
	num = 0
	fmt.Println("---------------------")
	fmt.Println("MAX INTERVAL FAILURE EXAMPLE")

	err = retry.NewExponentialBackoff(retry.ExponentialBackoff{
		MaxInterval: 10 * time.Second,
		MaxAttempts: 10,
	}).Retry(func() interface{} {
		stdErr = countNoReturn(num)
		return stdErr
	})
	if err != nil {
		fmt.Println(err)
	}
}

// EXAMPLE FUNCTION
func count(n int) (string, error) {
	if n != 5 {
		num++ // comment out to fail on maxInterval
		return "", errors.New("error number is not 5, number is " + strconv.Itoa(n))
	}
	return "counted to 5 successfully", nil
}

func countNoReturn(n int) error {
	if n != 5 {
		num++ // comment out to fail on maxInterval
		return errors.New("error number is not 5, number is " + strconv.Itoa(n))
	}
	return nil
}
