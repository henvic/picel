package logger

import (
	"log"
	"os"
)

var (
	Stdout = log.New(os.Stdout, "", log.LstdFlags)
	Stderr = log.New(os.Stderr, "", log.LstdFlags)
)
