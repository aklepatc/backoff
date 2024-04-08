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
		err := call()
		if err != nil && limit > 1 {
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
