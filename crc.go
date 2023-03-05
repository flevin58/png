package main

import "log"

type CRC struct {
	Table    [256]uint32
	Computed bool
	Value    uint32
}

func NewCRC() *CRC {
	var (
		crc CRC
		n   uint32
		c   uint32
	)

	crc = CRC{}
	crc.Computed = false
	for n = 0; n < 256; n++ {
		c = n
		for k := 0; k < 8; k++ {
			if c&1 != 0 {
				c = 0xedb88320 ^ (c >> 1)
			} else {
				c = c >> 1
			}
		}
		crc.Table[n] = c
	}
	crc.Computed = true
	crc.Value = 0xffffffff
	return &crc
}

func (c *CRC) Update(buf []byte) {
	if !c.Computed {
		log.Fatal("crc table not initialized. Must call NewCrc() function.")
	}

	for n := 0; n < len(buf); n++ {
		ch := uint32(buf[n])
		c.Value = c.Table[(c.Value^ch)&0xff] ^ (c.Value >> 8)
	}
}

func (c *CRC) Start(buf []byte) {
	c.Value = 0xffffffff
	c.Update(buf)
}

func (c *CRC) Get() uint32 {
	return c.Value ^ 0xffffffff
}
