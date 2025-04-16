package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

type STRUCT struct {
	X int16 `struc:"int16,little"`
	Y int16 `struc:"int16,little"`
	Z int16 `struc:"int16,little"`
	A bool  `struc:"bool"`
	B bool  `struc:"bool"`
}

func TestDecode(t *testing.T) {
	input, err := hex.DecodeString("7e00ff08a7025005aa020100aa81")
	assert.NoError(t, err)
	var msg STRUCT
	assert.NoError(t, DecodePacket(input, &msg))
	assert.Equal(t, STRUCT{
		X: 679,
		Y: 1360,
		Z: 682,
		A: true,
		B: false,
	}, msg)
}

func TestEncode(t *testing.T) {
	input := STRUCT{
		X: 679,
		Y: 1360,
		Z: 682,
		A: true,
		B: false,
	}
	packet, err := EncodePacket(&input)
	assert.NoError(t, err)
	assert.Equal(t, "7e00ff08a7025005aa020100aa81", hex.EncodeToString(packet))
}

func TestEncodeWithStartByte(t *testing.T) {
	input := STRUCT{
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

func TestDecode2(t *testing.T) {
	input, err := hex.DecodeString("7e00ff1210020402f201000101000000000000c819004781")
	assert.NoError(t, err)
	var msg message
	assert.NoError(t, DecodePacket(input, &msg))
	assert.Equal(t, message{
		X: 679,
		Y: 1360,
		Z: 682,
	}, msg)
}

//func TestReader(t *testing.T) {
//	input, err := hex.DecodeString("7e00ff08a7025005aa020100aa817e0000080002fe04810201006281")
//	assert.NoError(t, err)
//
//	reader := NewReader(bytes.NewReader(input))
//
//	var msg STRUCT
//	assert.NoError(t, reader.Read(&msg))
//	assert.Equal(t, STRUCT{
//		X: 679,
//		Y: 1360,
//		Z: 682,
//		A: true,
//		B: false,
//	}, msg)
//}
