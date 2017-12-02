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
	Run(context.Context, chan<- *CheckResult)
}

type commonChecker struct {
	Timeout time.Duration
}
