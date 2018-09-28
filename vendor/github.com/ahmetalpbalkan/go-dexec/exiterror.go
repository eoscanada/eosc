package dexec

import "fmt"

// ExitError reports an unsuccessful exit by a command.
type ExitError struct {
	// ExitCode holds the non-zero exit code of the container
	ExitCode int

	// Stderr holds the standard error output from the command
	// if it *Cmd executed through Output() and Cmd.Stderr was not
	// set.
	Stderr []byte
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("dexec: exit status: %d", e.ExitCode)
}
