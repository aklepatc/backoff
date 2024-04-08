package main

import (
	"context"
	"fmt"
	"time"
)

func expBackoff(
	ctx context.Context,
	call func() error,
	delay time.Duration,
	multiplier float64,
	limit int,
) error {
	done := ctx.Done()
	result := make(chan error, 1)
	go func() {
		if limit < 1 {
			result <- fmt.Errorf(`expected: limit>=1 got: limit=%d`, limit)
			return
		}
		err := call()
		switch {
		case err == nil || limit == 1:
			result <- err
		case multiplier < 1.0:
			result <- fmt.Errorf(`expected: multiplier>=1.0 got: multiplier=%v`, multiplier)
		default:
			timer := time.NewTimer(delay)
			fDelay := float64(delay)
			for step := 2; ; step++ {
				select {
				case <-done:
					result <- ctx.Err()
					timer.Stop()
					return
				case <-timer.C:
					err = call()
					if err == nil || step == limit {
						result <- err
						return
					}
					fDelay *= multiplier
					timer.Reset(time.Duration(fDelay))
				}
			}
		}
	}()
	select {
	case <-done:
		return ctx.Err()
	case err := <-result:
		return err
	}
}

func main() {
	fmt.Println("Hello, 世界")
	// TODO meaningful test logic
}
