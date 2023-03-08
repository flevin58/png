package pngfile

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/flevin58/png/errors"
)

// The PNG magic number at the beginning of all PNG files.
const (
	PNGSignature uint64 = 0x89504e470d0a1a0a
)

// The PNG file structure. Top level is just an array of chunks!
type PngType struct {
	MagicNumber uint64
	Chunks      []Chunk
}

// Read the PNG file and set the structure
func (p *PngType) Read(fileName string) {

	f, err := os.Open(fileName)
	errors.AbortOnError(err, "opening file")
	defer f.Close()

	err = binary.Read(f, binary.BigEndian, &p.MagicNumber)
	errors.AbortOnError(err, "reading png magic number")
	if p.MagicNumber != PNGSignature {
		errors.Error("file %s not recognized as a PNG file\n", fileName)
	}

	var ch Chunk = Chunk{}

	for ch.Type != IEND {
		ch.Read(f)
		p.Chunks = append(p.Chunks, ch)
	}
}

// Write the PNG struct to file
func (p *PngType) Write(fileName string) {

	log.Println("writing image file:", fileName)

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("error opening %v for save: %v\n", fileName, err)
	}
	defer f.Close()

	binary.Write(f, binary.BigEndian, p.MagicNumber)
	for _, ch := range p.Chunks {
		ch.Write(f)
	}
}

func (p *PngType) FindChunk(chunkType uint32) *Chunk {
	for _, ch := range p.Chunks {
		if ch.Type == chunkType {
			return &ch
		}
	}
	return nil
}
