package checkers

import "time"

func NewShellChecker(timeout time.Duration, command string) (Checker, error) {
	return &CommandChecker{
		commonChecker: commonChecker{Timeout: timeout},
		Exec:          []string{"sh", "-c", command},
	}, nil
}
