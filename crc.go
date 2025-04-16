package main

type PacketCRC struct {
	poly  byte
	table [256]byte
}

func NewPacketCRC(polynomial byte) *PacketCRC {
	crc := &PacketCRC{poly: polynomial}
	crc.generateTable()
	return crc
}

func (c *PacketCRC) generateTable() {
	for i := 0; i < 256; i++ {
		curr := byte(i)
		for j := 0; j < 8; j++ {
			if curr&0x80 != 0 {
				curr = (curr << 1) ^ c.poly
			} else {
				curr <<= 1
			}
		}
		c.table[i] = curr
	}
}

func (c *PacketCRC) Calculate(data []byte) byte {
	var crc byte = 0
	for _, b := range data {
		crc = c.table[crc^b]
	}
	return crc
}
