package errs

import (
	"fmt"
	"os"
)

// Error prints a notice to os.Stderr on what line an error has occurred as
// well as a brief message explaining the failure.
func Error(line uint, message string) {
	report(line, "", message)
}

func report(line uint, where string, message string) {
	if len(where) > 0 {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", line, where, message)
		return
	}

	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
}
