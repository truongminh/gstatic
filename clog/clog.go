package clog

import (
	"log"
	"os"
	"runtime/debug"
)

// Logger is the clog logging interface.
type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	V(int) bool
	SetV(level int)
}

var logger Logger = &stdlog{
	verbosity: 0,
}

// SetLogger set the clog logging implementation.
func SetLogger(l Logger) { logger = l }

var verbosity int

// V returns whether the current clog verbosity is above the specified level.
func V(level int) bool {
	if logger == nil {
		return false
	}
	return logger.V(level)
}

// SetV sets the clog verbosity level.
func SetV(level int) {
	if logger != nil {
		logger.SetV(level)
	}
}

// Infof logs information level messages.
func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Infof(format, args...)
	}
}

// Warningf logs warning level messages.
func Warningf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warningf(format, args...)
	}
}

// Errorf logs error level messages.
func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(format, args...)
	}
}

// Fatalf logs fatal messages and terminates the program.
func Fatalf(format string, args ...interface{}) {
	if logger != nil {
		logger.Fatalf(format, args...)
	}
}

// stdlog wraps the standard library logger.
type stdlog struct {
	verbosity int
}

func (stdlog) Infof(format string, args ...interface{})    { log.Printf(format, args...) }
func (stdlog) Warningf(format string, args ...interface{}) { log.Printf("WARN: "+format, args...) }
func (stdlog) Errorf(format string, args ...interface{})   { log.Printf("ERROR: "+format, args...) }
func (stdlog) Fatalf(format string, args ...interface{})   { log.Fatalf("FATAL: "+format, args...) }
func (s stdlog) V(level int) bool                          { return s.verbosity >= level }
func (s *stdlog) SetV(level int)                           { s.verbosity = level }

// prefixlog wraps the standard library logger.
type prefixlog struct {
	prefix    string
	verbosity int
}

func (s *prefixlog) Infof(format string, args ...interface{}) { log.Printf(s.prefix+format, args...) }
func (s *prefixlog) Warningf(format string, args ...interface{}) {
	log.Printf(s.prefix+"WARN: "+format, args...)
}
func (s *prefixlog) Errorf(format string, args ...interface{}) {
	if debugStackEnabled {
		debug.PrintStack()
	}
	log.Printf(s.prefix+"ERROR: "+format, args...)
}
func (s *prefixlog) Fatalf(format string, args ...interface{}) {
	log.Fatalf(s.prefix+"FATAL: "+format, args...)
}
func (s prefixlog) V(level int) bool { return s.verbosity >= level }
func (s *prefixlog) SetV(level int)  { s.verbosity = level }

// Prefix new prefix logger
func Prefix(prefix string) Logger {
	s := &prefixlog{verbosity: 0}
	s.prefix = "[" + prefix + "] "
	return s
}

var _, debugStackEnabled = os.LookupEnv("DEBUG_STACK_ENABLED")
