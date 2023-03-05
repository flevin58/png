// Convert struct to bytes and viceversa

package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"os"
)

/*
	var cIHDR []byte = []byte{'I', 'H', 'D', 'R'}
	var u32 uint32 = binary.LittleEndian.Uint32(cIHDR)
	fmt.Printf("cIHDR = %s %X\n", string(cIHDR), u32)
	os.Exit(0)
*/

// Known chunk types
const (
	IHDR uint32 = 0x52444849
	PLTE uint32 = 0x45544c50
	IDAT uint32 = 0x54414449
	IEND uint32 = 0x444e4549
)

type Chunk struct {
	Length uint32
	Type   uint32
	Data   []byte
	CRC    uint32
}

//---------------------[ IHDR ]----------------------

// ColourType values
const (
	Greyscale       byte = 0
	TrueColour      byte = 2
	IndexedColour   byte = 3
	GreyscaleAlpha  byte = 4
	TruecolourAlpha byte = 6
)

type ChIHDR struct {
	Width             uint32
	Height            uint32
	BitDepth          byte
	ColourType        byte
	CompressionMethod byte
	FilterMethod      byte
	InterlaceMethod   byte
}

func (ch *Chunk) String() string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, ch.Type)
	return fmt.Sprintf("Len: %d, Type: %s (0x%X), CRC: %X",
		ch.Length, string(bs), ch.Type, ch.CRC)
}

func (ch *Chunk) Read(f *os.File) {

	var err error

	// Read chunk from file
	err = binary.Read(f, binary.BigEndian, &ch.Length)
	if err != nil {
		log.Fatalf("reading length: %v", err)
	}
	err = binary.Read(f, binary.LittleEndian, &ch.Type)
	if err != nil {
		log.Fatalf("reading type: %v", err)
	}
	ch.Data = make([]byte, ch.Length)
	err = binary.Read(f, binary.LittleEndian, ch.Data)
	if err != nil {
		log.Fatalf("reading data: %v", err)
	}
	err = binary.Read(f, binary.BigEndian, &ch.CRC)
	if err != nil {
		log.Fatalf("reading crc: %v", err)
	}

	// Check CRC
	var (
		buf4 []byte = make([]byte, 4)
		crc  uint32 = 0
	)

	binary.LittleEndian.PutUint32(buf4, ch.Type)
	crc = crc32.Update(crc, crc32.IEEETable, buf4)
	crc = crc32.Update(crc, crc32.IEEETable, ch.Data)
	if ch.CRC != crc {
		log.Println(ch.String())
		log.Fatalf("crc error: expected %08X, calculated %08X", ch.CRC, crc)
	}
}

func (ch *Chunk) Write(f *os.File) {
	binary.Write(f, binary.BigEndian, ch.Length)
	binary.Write(f, binary.LittleEndian, ch.Type)
	binary.Write(f, binary.LittleEndian, ch.Data)
	binary.Write(f, binary.BigEndian, ch.CRC)
}
