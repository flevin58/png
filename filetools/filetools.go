package filetools

import (
	"os"

	"github.com/flevin58/png/errors"
)

func SafeBytesRead(f *os.File, buf []byte, msg string) {
	n, err := f.Read(buf)
	if n != len(buf) || err != nil {
		errors.Error("%s: %v\n", msg, err)
	}
}
