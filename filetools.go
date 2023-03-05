package main

import (
	"log"
	"os"
)

func SafeBytesRead(f *os.File, buf []byte, msg string) {
	n, err := f.Read(buf)
	if n != len(buf) || err != nil {
		log.Fatalf("error %s: %v\n", msg, err)
	}
}
