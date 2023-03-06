package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func doCopy(inf interface{}) {
	args := inf.(*CmdCopy)
	image.Read(args.Ifile)
	image.Write(args.Ofile)
}

func doDumpHeader() {
	var header ChIHDR = ChIHDR{}

	ch := image.FindChunk(IHDR)
	if ch == nil {
		Error("chunk type '%s' not found", ch.StrType())
	}

	buffer := &BytesReadWriteSeeker{
		data: ch.Data,
		pos:  0,
	}
	buffer.Seek(0, io.SeekStart)
	binary.Read(buffer, binary.BigEndian, &header.Width)
	binary.Read(buffer, binary.BigEndian, &header.Height)

	fmt.Printf("Size..............: %d x %d\n", header.Width, header.Height)
	fmt.Printf("Bit depth.........: %d\n", header.BitDepth)
	fmt.Printf("Colour type.......: %d\n", header.ColourType)
	fmt.Printf("Compression method: %d\n", header.CompressionMethod)
	fmt.Printf("Filter method.....: %d\n", header.FilterMethod)
	fmt.Printf("Interlace method..: %d\n", header.InterlaceMethod)
}

func doDumpPalette() {
	Abort("TODO: dump PLTE chunk\n")
}

func doDumpAll() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Type", "Length (bytes)", "CRC"})
	for i, ch := range image.Chunks {
		t.AppendRow([]interface{}{
			i,
			ch.StrType(),
			ch.Length,
			fmt.Sprintf("%08X", ch.CRC),
		})
		t.AppendSeparator()
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func doDump(inf interface{}) {
	args := inf.(*CmdDump)
	image.Read(args.Ifile)

	switch args.Chunk {
	case "":
		doDumpAll()
	case "IHDR":
		doDumpHeader()
	case "PLTE":
		doDumpPalette()
	default:
		Error("unrecognized or unsupported chunk: %v\n", args.Chunk)
	}
}
