package serialtransfer

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	X int16 `struc:"int16,little"`
	Y int16 `struc:"int16,little"`
	Z int16 `struc:"int16,little"`
	A bool  `struc:"bool"`
	B bool  `struc:"bool"`
}

func TestDecode(t *testing.T) {
	input, err := hex.DecodeString("7e00ff08a7025005aa020100aa81")
	assert.NoError(t, err)
	var msg testMessage
	assert.NoError(t, DecodePacket(input, &msg))
	assert.Equal(t, testMessage{
		X: 679,
		Y: 1360,
		Z: 682,
		A: true,
		B: false,
	}, msg)
}

func TestEncodeWithStartByte(t *testing.T) {
	input := testMessage{
		X: 638,
		Y: 1278,
		Z: 641,
		A: true,
		B: false,
	}
	packet, err := EncodePacket(&input)
	assert.NoError(t, err)
	assert.Equal(t, "7e0000080002fe04810201006281", hex.EncodeToString(packet))
}

func TestDecodePayloadEncodedWithLargerStruct(t *testing.T) {
	input, err := hex.DecodeString("7e00ff1210020402f201000101000000000000c819004781")
	assert.NoError(t, err)
	var msg testMessage
	assert.NoError(t, DecodePacket(input, &msg))
	assert.Equal(t, testMessage{
		X: 528,
		Y: 516,
		Z: 498,
		A: false,
		B: true,
	}, msg)
}

func TestReader(t *testing.T) {
	input, err := hex.DecodeString(
		"00ff" +
			"7e00ff08a7025005aa020100aa81" +
			"7e0000080002fe04810201006281")
	assert.NoError(t, err)

	reader := NewDecoder(bytes.NewReader(input))

	var msgs []testMessage
	for {
		var msg testMessage
		err := reader.Decode(&msg)
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		msgs = append(msgs, msg)
	}
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, []testMessage{
		{
			X: 679,
			Y: 1360,
			Z: 682,
			A: true,
			B: false,
		},
		{
			X: 638,
			Y: 1278,
			Z: 641,
			A: true,
			B: false,
		},
	}, msgs)
}
