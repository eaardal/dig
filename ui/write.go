package ui

import (
	"fmt"
	"os"
)

func Write(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format+"\n", args...)
}
