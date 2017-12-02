package checkers

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/juju/errors"

	log "github.com/sirupsen/logrus"
)

type CommandChecker struct {
	commonChecker

	Exec []string
}

func (cc *CommandChecker) StrCommand() string {
	var buf bytes.Buffer

	for _, v := range cc.Exec {
		if strings.ContainsRune(v, ' ') {
			buf.WriteRune('"')
			buf.WriteString(v)
			buf.WriteRune('"')
		} else {
			buf.WriteString(v)
		}
		buf.WriteRune(' ')
	}

	return strings.TrimSpace(buf.String())
}

func (cc *CommandChecker) Run(ctx context.Context, results chan<- *CheckResult) {
	cmd := cc.StrCommand()
	newCtx, cancel := context.WithTimeout(ctx, cc.Timeout)
	defer cancel()

	log.WithFields(log.Fields{
		"command": cmd,
		"timeout": cc.Timeout,
	}).Debug("Run command")

	result := CheckResult{Ok: true, Producer: "command: " + cmd}
	if err := exec.CommandContext(newCtx, cc.Exec[0], cc.Exec[1:]...).Run(); err != nil {
		result.Ok = false
		result.Error = errors.Annotate(err, "Command "+cmd+" failed")
	}
	results <- &result
}

func NewCommandChecker(timeout time.Duration, command string) (Checker, error) {
	parsedCommand, err := shlex.Split(command)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &CommandChecker{
		commonChecker: commonChecker{Timeout: timeout},
		Exec:          parsedCommand,
	}, nil
}
