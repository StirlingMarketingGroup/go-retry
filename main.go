package retry

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// Min returns the min of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the max of two ints
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RandInt returns a random int between (including) the min and max
func RandInt(min, max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(nBig.Int64()) + min, nil
}

// Retry tries to execute a function with "retries"+1 times until it's succeeded
func Retry(main func() error, retries int, afterTryFailure func(error) error, beforeRetry func() error) error {
	var mainErr error

	if main == nil {
		return fmt.Errorf("the main function to try can't be nil")
	}

	for i := 0; i <= retries; i++ {
		if i != 0 && beforeRetry != nil {
			err := beforeRetry()
			if err != nil {
				return fmt.Errorf("retry before function: %s", err)
			}
		}

		mainErr = main()
		if mainErr == nil {
			break
		}

		if afterTryFailure != nil {
			err := afterTryFailure(mainErr)
			if err != nil {
				return fmt.Errorf("retry after function: %s", err)
			}
		}
		if i == retries {
			break
		}
		n, err := RandInt(-2, 5)
		if err != nil {
			return fmt.Errorf("retry rand: %s", err)
		}
		wait := Min(15000, Max(i*750, i+n*100))
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}

	return mainErr
}
