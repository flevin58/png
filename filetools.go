package main

import (
	"os"
)

func SafeBytesRead(f *os.File, buf []byte, msg string) {
	n, err := f.Read(buf)
	if n != len(buf) || err != nil {
		Error("%s: %v\n", msg, err)
	}
}
