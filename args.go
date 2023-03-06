package main

import (
	"github.com/alexflint/go-arg"
)

type InputFileArg struct {
	Ifile string `arg:"positional, required" help:"The name of the source file"`
}

type OutputFileArg struct {
	Ofile string `arg:"positional, required" help:"The name of the output file"`
}

type CmdCopy struct {
	InputFileArg
	OutputFileArg
}

type CmdDump struct {
	InputFileArg
	Chunk string `arg:"-c,--chunk" default:"" help:"dump the chunk with that name"`
}

type Args struct {
	Copy *CmdCopy `arg:"subcommand"`
	Dump *CmdDump `arg:"subcommand"`
}

var args Args

// Parses the command line and returns the subcommand name and the pointer to its arg structure
func ParseCommandLine() (string, interface{}) {
	p := arg.MustParse(&args)
	return p.SubcommandNames()[0], p.Subcommand()
}
