package serialtransfer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lunixbochs/struc"
)

const (
	StartByte = 0x7E
	StopByte  = 0x81
)

func DecodePacket(data []byte, out interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(out)
}

type Decoder struct {
	crcChecker *PacketCRC
	next       io.Reader

	state fsm

	packetID     byte
	overheadByte byte
	packetLen    int

	payload      []byte
	payloadIndex int
}

type fsm int

const (
	findStartByte fsm = iota
	findIdByte
	findOverheadByte
	findPayloadLen
	findPayload
	findCrc
	findEndByte
)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		crcChecker:   NewPacketCRC(0x9B),
		next:         r,
		payload:      make([]byte, 255),
		payloadIndex: 0,
		state:        findStartByte,
	}
}

func (r *Decoder) Decode(data interface{}) error {
	for {
		b := make([]byte, 1)
		_, err := r.next.Read(b)
		if err != nil {
			return err
		}

		receivedByte := b[0]
		switch r.state {
		case findStartByte:
			if receivedByte == StartByte {
				r.state = findIdByte
			}
		case findIdByte:
			r.packetID = receivedByte
			r.state = findOverheadByte
		case findOverheadByte:
			r.overheadByte = receivedByte
			r.state = findPayloadLen
		case findPayloadLen:
			r.packetLen = int(receivedByte)
			r.state = findPayload
			r.payloadIndex = 0
		case findPayload:
			if r.payloadIndex >= r.packetLen {
				return fmt.Errorf("payload index out of range")
			}
			r.payload[r.payloadIndex] = receivedByte
			r.payloadIndex++
			if r.payloadIndex == r.packetLen {
				r.state = findCrc
			}
		case findCrc:
			calculatedCRC := r.crcChecker.Calculate(r.payload[:r.payloadIndex])
			if calculatedCRC != receivedByte {
				r.state = findStartByte
				return fmt.Errorf("crc check failed: expected 0x%X, got 0x%X", receivedByte, calculatedCRC)
			} else {
				r.state = findEndByte
			}
		case findEndByte:
			cobsPayload := make([]byte, r.packetLen)
			copy(cobsPayload, r.payload[:r.payloadIndex])
			cobsDecode(cobsPayload, r.overheadByte)

			if err := struc.Unpack(bytes.NewReader(cobsPayload), data); err != nil {
				return fmt.Errorf("failed to unpack: %v", err)
			}
			r.state = findStartByte
			return nil
		}
	}
}

func cobsDecode(arr []byte, recOverheadByte byte) {
	var testIndex = recOverheadByte
	var delta byte = 0

	for int(testIndex) < len(arr) && arr[testIndex] != 0 {
		delta = arr[testIndex]
		arr[testIndex] = StartByte
		testIndex += delta
	}
	if int(testIndex) < len(arr) {
		arr[testIndex] = StartByte
	}
}
