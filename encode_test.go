package serialtransfer

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	input := testMessage{
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
