/*
Package logger provides logging for picel.
*/
package logger

import (
	"log"
	"os"
)

var (
	// Stdout is the standard output logger
	Stdout = log.New(os.Stdout, "", log.LstdFlags)

	// Stderr is the standard error output logger
	Stderr = log.New(os.Stderr, "", log.LstdFlags)
)
