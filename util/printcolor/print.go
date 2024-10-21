package printcolor

import (
	"github.com/fatih/color"
)

// Error prints an error message in red
func Error(msg string, args ...interface{}) {
	color.Red(msg, args...)
}

// Success prints a success message in green
func Success(msg string, args ...interface{}) {
	color.Green(msg, args...)
}

// Info prints an info message in blue
func Info(msg string, args ...interface{}) {
	color.Blue(msg, args...)
}

// Warning prints a warning message in yellow
func Warning(msg string, args ...interface{}) {
	color.Yellow(msg, args...)
}

// Print prints a message in white
func Print(msg string, args ...interface{}) {
	color.White(msg, args...)
}
