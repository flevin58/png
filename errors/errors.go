package errors

import (
	"fmt"
	"os"
)

func Warning(format string, args ...any) {
	fmt.Printf("WARNING: "+format, args...)
}

func Error(format string, args ...any) {
	fmt.Printf("ERROR: "+format, args...)
	os.Exit(1)
}

func Abort(format string, args ...any) {
	fmt.Printf("ABORT: "+format, args...)
	os.Exit(1)
}

func AbortOnError(err error, format string, args ...any) {
	if err != nil {
		fmt.Printf("ERROR: "+format+":%v\n", args...)
		os.Exit(1)
	}
}
