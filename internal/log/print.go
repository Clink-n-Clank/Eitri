package log

import (
	"fmt"
	"os"
)

// PrintError will display CMD error
func PrintError(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[1;31m ERROR: %s \033[m\n", msg)
}

// PrintInfo will display Info in CMD
func PrintInfo(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "\033[1;32m INFO: %s \033[m\n", msg)
}
