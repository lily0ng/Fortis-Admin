package logging

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger struct {
	out     io.Writer
	errOut  io.Writer
	level   Level
	quiet   bool
	verbose bool
}

func New(out, errOut io.Writer, level Level, quiet, verbose bool) *Logger {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &Logger{out: out, errOut: errOut, level: level, quiet: quiet, verbose: verbose}
}

func (l *Logger) SetLevel(level Level) { l.level = level }

func (l *Logger) Debugf(format string, args ...any) {
	if l.quiet || !l.verbose {
		return
	}
	if l.level < LevelDebug {
		return
	}
	l.printf(l.out, "DEBUG", format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	if l.quiet {
		return
	}
	if l.level < LevelInfo {
		return
	}
	l.printf(l.out, "INFO", format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	if l.quiet {
		return
	}
	if l.level < LevelWarn {
		return
	}
	l.printf(l.errOut, "WARN", format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.level < LevelError {
		return
	}
	l.printf(l.errOut, "ERROR", format, args...)
}

func (l *Logger) Println(msg string) {
	if l.quiet {
		return
	}
	fmt.Fprintln(l.out, msg)
}

func (l *Logger) printf(w io.Writer, tag, format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s [%s] %s\n", ts, tag, fmt.Sprintf(format, args...))
}
