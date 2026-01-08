package ui

import (
	"fmt"
	"io"
)

type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

func Banner(w io.Writer, title string) {
	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "║ %-56s ║\n", title)
	fmt.Fprintln(w, "║         System Administration & Security Toolkit         ║")
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)
}

func colorize(enabled bool, code, s string) string {
	if !enabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func Red(enabled bool, s string) string    { return colorize(enabled, "31", s) }
func Green(enabled bool, s string) string  { return colorize(enabled, "32", s) }
func Yellow(enabled bool, s string) string { return colorize(enabled, "33", s) }
func Blue(enabled bool, s string) string   { return colorize(enabled, "34", s) }
