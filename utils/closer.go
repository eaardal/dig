package utils

import (
	"fmt"
)

func CloseOrPanic(closeFn func() error) {
	if err := closeFn(); err != nil {
		panic(fmt.Errorf("could not close: %w", err))
	}
}
