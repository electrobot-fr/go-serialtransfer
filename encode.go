package serialtransfer

import (
	"bytes"

	"github.com/lunixbochs/struc"
)

func EncodePacket(s interface{}) ([]byte, error) {
	return NewEncoder().Encode(s)
}

type Encoder struct {
	crcChecker *PacketCRC
}

func NewEncoder() *Encoder {
	return &Encoder{
		crcChecker: NewPacketCRC(0x9B),
	}
}

func (e *Encoder) Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := struc.Pack(&buf, data); err != nil {
		return nil, err
	}
	packed := buf.Bytes()
	overheadByte := findLast(packed)
	cobsEncode(packed, overheadByte)
	calculatedCRC := e.crcChecker.Calculate(packed)
	return append(append(append(
		[]byte{StartByte, 0x00, byte(overheadByte), byte(buf.Len())}),
		packed...),
		[]byte{calculatedCRC, StopByte}...), nil
}

func cobsEncode(arr []byte, overheadByte int) {
	length := len(arr)
	refByte := overheadByte

	if refByte != -1 {
		for i := length - 1; i >= 0; i-- {
			if arr[i] == StartByte {
				arr[i] = byte(refByte - i)
				refByte = i
			}
		}
	}
}

func findLast(arr []byte) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == StartByte {
			return i
		}
	}
	return -1
}
