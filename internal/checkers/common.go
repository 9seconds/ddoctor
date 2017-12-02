package checkers

import (
	"context"
	"time"
)

type CheckResult struct {
	Ok       bool
	Error    error
	Producer string
}

type Checker interface {
	Run(ctx context.Context, results chan<- *CheckResult)
}

type commonChecker struct {
	Timeout time.Duration
}
