package main

import (
	"fmt"
)

func Info(format string, args ...any) {
	fmt.Print("\x1b[1;7;36m info  \x1b[0m ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func Warn(format string, args ...any) {
	fmt.Print("\x1b[1;7;33m warn  \x1b[0m ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func Error(format string, args ...any) {
	fmt.Print("\x1b[1;7;31m error \x1b[0m ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func Trace(format string, args ...any) {
	fmt.Print("\x1b[1;7;90m debug \x1b[0m ")
	fmt.Printf(format, args...)
	fmt.Println()
}
