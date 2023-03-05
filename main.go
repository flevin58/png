//
// Based on: https://www.w3.org/TR/png/
//

package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

type Args struct {
	Ifile string `arg:"positional, required" help:"The source file"`
	Ofile string `arg:"positional, required" help:"The output file"`
}

var (
	args  Args
	image PngFile
	err   error
	flog  *os.File
)

func main() {

	// Initialize log file
	flog, err = os.OpenFile("pngfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer flog.Close()
	log.SetOutput(flog)

	// Parse the args and copy the file
	arg.MustParse(&args)
	image.Read(args.Ifile)
	image.Write(args.Ofile)
}
