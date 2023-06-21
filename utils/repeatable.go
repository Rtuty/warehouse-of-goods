package utils

import "time"

func DoWithTries(fn func() error, attemptps int, delay time.Duration) (err error) {
	for attemptps > 0 {
		if err := fn(); err != nil {
			time.Sleep(delay)
			attemptps--

			continue
		}

		return nil
	}
	return
}
