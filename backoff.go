package main

import (
	"context"
	"fmt"
	"time"
)

func expBackoff(
	ctx context.Context,
	callback func() error,
	delay time.Duration,
	multiplier int64,
	limit int,
) error {
	done := ctx.Done()
	select {
	case <-done:
		return ctx.Err()
	default:
	}
	result := make(chan error, 1)
	go func() {
		err := callback()
		if err != nil && limit > 1 {
			delayNS := delay.Nanoseconds()
			timer := time.NewTimer(time.Duration(delayNS))
			for step := 1; ; {
				select {
				case <-done:
					result <- ctx.Err()
					timer.Stop()
					return
				case <-timer.C:
					err = callback()
					step++
					if err == nil || step == limit {
						result <- err
						return
					}
					delayNS *= multiplier
					timer.Reset(time.Duration(delayNS))
				}
			}
		}
		result <- err
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
