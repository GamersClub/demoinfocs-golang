package sendtablescs2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFindDecoder_Fixed8ReadsRawBytes verifies that "fixed8" fields are read as 8 raw bits. Read as
// a varint, 0xFF would also consume the next byte (CS2 update of 2026-11-11, m_recentKillQueue).
func TestFindDecoder_Fixed8ReadsRawBytes(t *testing.T) {
	t.Parallel()

	f := field{varName: "m_recentKillQueue", varType: "uint8[8]", encoder: "fixed8", fieldType: newFieldType("uint8")}
	dec := findDecoder(&f)
	r := newReader([]byte{0xFF, 0x01, 0x80})

	assert.Equal(t, uint64(0xFF), dec(r))
	assert.Equal(t, uint64(0x01), dec(r))
	assert.Equal(t, uint64(0x80), dec(r))
	assert.Zero(t, r.remBytes())
}

// TestFindDecoder_Fixed16And32 verifies the wider fixedN encoders read little-endian raw bits.
func TestFindDecoder_Fixed16And32(t *testing.T) {
	t.Parallel()

	f16 := field{encoder: "fixed16", fieldType: newFieldType("uint16")}
	assert.Equal(t, uint64(0xBEEF), findDecoder(&f16)(newReader([]byte{0xEF, 0xBE})))

	f32 := field{encoder: "fixed32", fieldType: newFieldType("uint32")}
	assert.Equal(t, uint64(0xDEADBEEF), findDecoder(&f32)(newReader([]byte{0xEF, 0xBE, 0xAD, 0xDE})))
}

// TestFindDecoder_Uint8WithoutEncoderStaysVarint guards the default: no encoder is still a varint.
func TestFindDecoder_Uint8WithoutEncoderStaysVarint(t *testing.T) {
	t.Parallel()

	f := field{varName: "m_x", varType: "uint8", fieldType: newFieldType("uint8")}

	assert.Equal(t, uint64(300), findDecoder(&f)(newReader([]byte{0xAC, 0x02})))
}
