package main

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
)

// The PNG magic number at the beginning of all PNG files.
const (
	PNGSignature uint64 = 0x89504e470d0a1a0a
)

// The PNG file structure. Top level is just an array of chunks!
type PngFile struct {
	MagicNumber [8]byte
	Chunks      []Chunk
}

// Read the PNG file and set the structure
func (p *PngFile) Read(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}
	defer f.Close()

	log.Println(strings.Repeat("-", 80))
	log.Println("reading image file:", fileName)
	buf8 := make([]byte, 8)
	SafeBytesRead(f, buf8, "error reading magic number")
	magicNum := binary.BigEndian.Uint64(buf8)
	if magicNum != PNGSignature {
		log.Fatalf("file %s not recognized as a PNG file\n", fileName)
	}

	var ch Chunk = Chunk{}

	for ch.Type != IEND {
		ch.Read(f)
		p.Chunks = append(p.Chunks, ch)
		log.Println(ch.String())
	}
}

// Write the PNG struct to file
func (p *PngFile) Write(fileName string) {
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
