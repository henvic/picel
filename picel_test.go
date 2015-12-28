package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/henvic/picel/logger"
)

type ExistsDependencyProvider struct {
	cmd  string
	find bool
}

func TestVersion(t *testing.T) {
	flagVersion = true
	main()
}

func TestExistsDependency(t *testing.T) {
	for _, c := range existsDependencyCases {
		exists := existsDependency(c.cmd)

		if exists != c.find {
			t.Errorf("existsDependency(%v) should return %v", c.cmd, c.find)
		}
	}
}

type CheckMissingDependenciesProvider struct {
	cmds      []string
	allExists bool
}

func TestCheckMissingDependencies(t *testing.T) {
	for _, c := range CheckMissingDependencies {
		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout := logger.Stdout
		defaultStderr := logger.Stderr
		logger.Stdout = log.New(&StdoutMock, "", log.LstdFlags)
		logger.Stderr = log.New(&StderrMock, "", log.LstdFlags)
		checkMissingDependencies(c.cmds...)
		logger.Stdout = defaultStdout
		logger.Stderr = defaultStderr

		if StdoutMock.String() != "" {
			t.Errorf("checkMissingDependencies(%v) stdout should be empty", c.cmds)
		}
	}
}
