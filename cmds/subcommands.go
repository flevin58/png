package cmds

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/flevin58/png/bytesrw"
	"github.com/flevin58/png/errors"
	"github.com/flevin58/png/pngfile"
	"github.com/jedib0t/go-pretty/v6/table"
)

var image pngfile.PngType

func DoCopy(inf interface{}) {
	args := inf.(*CmdCopy)
	image.Read(args.Ifile)
	image.Write(args.Ofile)
}

func doBinaryDump(chType string) {

	var (
		ascii  string
		char   byte
		i, j   uint32
		offset uint32
	)

	ch := image.FindChunk(pngfile.ChunkStringToUint32(chType))
	if ch == nil {
		errors.Error("chunk type '%s' not found\n", chType)
	}

	fmt.Printf("Dumping chunk '%s' of length %d\n\n", chType, ch.Length)
	for i = 0; i < ch.Length; i += 16 {
		ascii = ""
		fmt.Printf("%08X:  ", i)
		for j = 0; j < 16; j++ {
			offset = i + j
			if offset >= ch.Length {
				fmt.Print("-- ")
				ascii += "."
				continue
			}

			char = ch.Data[offset]
			fmt.Printf("%02X ", char)
			if char <= 32 || char >= 127 {
				ascii += "."
			} else {
				ascii += string(char)
			}
		}
		fmt.Println("  ", ascii)
	}
}

func doDumpHeader() {
	var (
		header   pngfile.ChIHDR
		cmString string
		fmString string
		imString string
	)

	header = pngfile.ChIHDR{}

	ch := image.FindChunk(pngfile.IHDR)
	if ch == nil {
		errors.Error("chunk type '%s' not found", ch.StrType())
	}

	buffer := bytesrw.NewBytesReadWriteSeeker(ch.Data)
	buffer.Seek(0, io.SeekStart)
	binary.Read(buffer, binary.BigEndian, &header.Width)
	binary.Read(buffer, binary.BigEndian, &header.Height)
	header.BitDepth, _ = buffer.ReadByte()
	header.ColourType, _ = buffer.ReadByte()
	header.CompressionMethod, _ = buffer.ReadByte()
	if header.CompressionMethod == 0 {
		cmString = "Deflate/inflate compression with a 32K sliding window."
	} else {
		cmString = "Unknown"
	}
	header.FilterMethod, _ = buffer.ReadByte()
	if header.CompressionMethod == 0 {
		fmString = "Adaptive filtering with five basic filter types."
	} else {
		fmString = "Unknown"
	}
	header.InterlaceMethod, _ = buffer.ReadByte()
	switch header.InterlaceMethod {
	case 0:
		imString = "No interlace."
	case 1:
		imString = "Adam7 interlace."
	default:
		imString = "Unknown"
	}

	fmt.Printf("Size..............: %d x %d\n", header.Width, header.Height)
	fmt.Printf("Bit depth.........: %d (Values range from 0 to %d)\n", header.BitDepth, uint(math.Pow(2, float64(header.BitDepth)))-1)
	fmt.Printf("Colour type.......: %d (%s)\n", header.ColourType, pngfile.ColorTypes[header.ColourType])
	fmt.Printf("Compression method: %d (%s)\n", header.CompressionMethod, cmString)
	fmt.Printf("Filter method.....: %d (%s)\n", header.FilterMethod, fmString)
	fmt.Printf("Interlace method..: %d (%s)\n", header.InterlaceMethod, imString)
}

func doDumpPalette() {
	ch := image.FindChunk(pngfile.PLTE)
	if ch == nil {
		errors.Error("chunk type '%s' not found", ch.StrType())
	}

	if ch.Length%3 != 0 {
		errors.Error("palette length %d is not multiple of 3 (r,g,b)\n", ch.Length)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Color Palette")
	t.AppendHeader(table.Row{"#", "R", "G", "B", "Hex"})
	for i, j := 0, 0; i < int(ch.Length); i += 3 {
		t.AppendRow([]interface{}{
			j,
			ch.Data[i],
			ch.Data[i+1],
			ch.Data[i+2],
			fmt.Sprintf("#%02X%02X%02X", ch.Data[i], ch.Data[i+1], ch.Data[i+2]),
		})
		t.AppendSeparator()
		j++
	}

	t.SetStyle(table.StyleLight)
	t.Render()
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

func DoDump(inf interface{}) {
	args := inf.(*CmdDump)
	image.Read(args.Ifile)

	switch args.Chunk {
	case "":
		doDumpAll()
	case "IHDR":
		doDumpHeader()
	case "PLTE":
		doBinaryDump(args.Chunk)
		doDumpPalette()
	default:
		doBinaryDump(args.Chunk)
	}
}
