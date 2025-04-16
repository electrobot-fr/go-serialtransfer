package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lunixbochs/struc"
	"github.com/sirupsen/logrus"
)

const (
	StartByte = 0x7E
	StopByte  = 0x81
)

func EncodePacket(s interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := struc.Pack(&buf, s); err != nil {
		return nil, err
	}
	packed := buf.Bytes()
	overheadByte := findLast(packed)
	cobsEncode(packed, overheadByte)
	crcChecker := NewPacketCRC(0x9B)
	calculatedCRC := crcChecker.Calculate(packed)
	return append(append(append(
		[]byte{StartByte, 0x00, byte(overheadByte), byte(buf.Len())}),
		packed...),
		[]byte{calculatedCRC, StopByte}...), nil
}

func DecodePacket(data []byte, out interface{}) error {
	if data[0] != StartByte || data[len(data)-1] != StopByte {
		return errors.New("invalid start/stop byte")
	}

	packetID := data[1]
	recOverheadByte := data[2]
	payloadLen := int(data[3])

	if len(data) < 5+payloadLen {
		return errors.New("payload length mismatch")
	}

	if payloadLen == 0 {
		return errors.New("payload is empty")
	}

	payloadStart := 4
	receivedCRC := data[payloadStart+payloadLen]
	cobsPayload := make([]byte, payloadLen)
	copy(cobsPayload, data[payloadStart:payloadStart+payloadLen])

	crcChecker := NewPacketCRC(0x9B)
	calculatedCRC := crcChecker.Calculate(cobsPayload)

	logrus.Debugf("Packet ID: %d, Overhead: %d, CRC: 0x%X, Calculated CRC: 0x%X\n", packetID, recOverheadByte, receivedCRC, calculatedCRC)

	if calculatedCRC != receivedCRC {
		return fmt.Errorf("crc check failed: expected 0x%X, got 0x%X", receivedCRC, calculatedCRC)
	}

	cobsDecode(cobsPayload, recOverheadByte)

	if err := struc.Unpack(bytes.NewReader(cobsPayload), out); err != nil {
		return fmt.Errorf("failed to unpack: %v", err)
	}
	return nil
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

//type Reader struct {
//	next io.Reader
//	buf  []byte
//	pos  int
//}
//
//func NewReader(r io.Reader) *Reader {
//	return &Reader{
//		next: r,
//		buf:  make([]byte, 256),
//		pos:  0,
//	}
//}
//
//func (r *Reader) Read(data interface{}) error {
//	b := make([]byte, 1)
//	_, err := r.next.Read(b)
//	if err != nil {
//		return err
//	}
//	if b[0] == StartByte {
//		DecodePacket(r.buf[:r.pos], data)
//	}
//	return nil
//}
