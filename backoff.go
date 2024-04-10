package main

import (
	"context"
	"fmt"
	"time"
)

func expBackoff(
	ctx context.Context,
	call func() error,
	limit int,
	delay time.Duration,
	multiplier float64,
) error {
	result := make(chan error, 1)
	done := ctx.Done()
	go func() {
		if limit < 1 || multiplier < 1.0 || delay == 0 {
			result <- fmt.Errorf(
				`want: limit>=1 got: limit=%d want: multiplier>=1.0 got: multiplier=%v want: delay>0 got: delay=%v`,
				limit, multiplier, delay,
			)
			return
		}
		select {
		case <-done:
			return
		default:
		}
		err := call()
		if err == nil || limit == 1 {
			result <- err
			return
		}
		timer := time.NewTimer(delay)
		for step, fDelay := 2, float64(delay); ; step++ {
			select {
			case <-done:
				timer.Stop()
				return
			case <-timer.C:
				err = call()
				if err == nil || step == limit {
					result <- err
					return
				}
			}
			fDelay *= multiplier
			timer.Reset(time.Duration(fDelay))
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
