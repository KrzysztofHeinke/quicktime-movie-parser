package atoms

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FixedPointToFloat32 converts a fixed-point 16.16 value to a float32
func FixedPointToFloat32(value uint32) float64 {
	return float64(value) / 65536.0
}

// TestParseStsdAtom tests the ParseStsdAtom function
func TestParseStsdAtom(t *testing.T) {
	// Prepare example data for the 'stsd' atom
	stsdData := []byte{
		0x00,             // Version
		0x00, 0x00, 0x00, // Flags
		0x00, 0x00, 0x00, 0x01, // EntryCount (1)
		// SampleEntry
		0x00, 0x00, 0x00, 0x58, // Size (88 bytes for this entry)
		'm', 'p', '4', 'a', // Type (mp4a)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x01, // Data reference index
		// Data (adjusted based on expected structure)
		0x00, 0x00, 0x00, 0x00, // Placeholder for additional data
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xBB, 0x80, 0x00, 0x00, // Sample rate 48000 (fixed point 16.16 format) at position 16:20
		0x00, 0x00, 0x00, 0x00, // Additional padding
	}

	reader := bytes.NewReader(stsdData)
	stsdAtom, err := ParseStsdAtom(reader)
	assert.NoError(t, err, "Expected no error parsing stsd atom")
	assert.Equal(t, uint32(1), stsdAtom.EntryCount, "Expected entry count to be 1")
	assert.Equal(t, "mp4a", string(stsdAtom.SampleEntries[0].Type[:]), "Expected entry type to be 'mp4a'")
}

// TestGetSampleRates tests the GetSampleRates function
func TestGetSampleRates(t *testing.T) {
	// Prepare a mock AtomStsd with a sample entry for 'mp4a'
	stsd := &AtomStsd{
		Version:    0,
		Flags:      [3]byte{0, 0, 0},
		EntryCount: 1,
		SampleEntries: []SampleEntry{
			{
				Size:     [4]byte{0x00, 0x00, 0x00, 0x58},
				Type:     [4]byte{'m', 'p', '4', 'a'},
				Reserved: [6]byte{0, 0, 0, 0, 0, 0},
				RefIndex: 1,
				Data: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0xBB, 0x80, 0x00, 0x00, // Sample rate 48000 (fixed point 16.16 format) at position 16:20
				},
			},
		},
	}

	sampleRates, err := GetSampleRates(stsd)
	assert.NoError(t, err, "Expected no error getting sample rates")
	assert.NotEmpty(t, sampleRates, "Expected sample rates to be found")
	assert.Contains(t, sampleRates, "mp4a", "Expected 'mp4a' codec to be present")
	assert.Equal(t, 48000.0, FixedPointToFloat32(uint32(sampleRates["mp4a"][0])), "Expected sample rate to be 48000.0 Hz")
}

// TestGetSampleRates_MultipleEntries tests GetSampleRates with multiple entries
func TestGetSampleRates_MultipleEntries(t *testing.T) {
	// Prepare a mock AtomStsd with multiple sample entries
	stsd := &AtomStsd{
		Version:    0,
		Flags:      [3]byte{0, 0, 0},
		EntryCount: 2,
		SampleEntries: []SampleEntry{
			{
				Size:     [4]byte{0x00, 0x00, 0x00, 0x58},
				Type:     [4]byte{'m', 'p', '4', 'a'},
				Reserved: [6]byte{0, 0, 0, 0, 0, 0},
				RefIndex: 1,
				Data: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0xBB, 0x80, 0x00, 0x00, // Sample rate 48000 (fixed point 16.16 format) at position 16:20
				},
			},
			{
				Size:     [4]byte{0x00, 0x00, 0x00, 0x58},
				Type:     [4]byte{'a', 'l', 'a', 'c'},
				Reserved: [6]byte{0, 0, 0, 0, 0, 0},
				RefIndex: 1,
				Data: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0xAC, 0x44, 0x00, 0x00, // Sample rate 44100 (fixed point 16.16 format) at position 16:20
				},
			},
		},
	}

	sampleRates, err := GetSampleRates(stsd)
	assert.NoError(t, err, "Expected no error getting sample rates")
	assert.NotEmpty(t, sampleRates, "Expected sample rates to be found")

	assert.Contains(t, sampleRates, "mp4a", "Expected 'mp4a' codec to be present")
	assert.Equal(t, 48000.0, FixedPointToFloat32(uint32(sampleRates["mp4a"][0])), "Expected sample rate to be 48000.0 Hz")

	assert.Contains(t, sampleRates, "alac", "Expected 'alac' codec to be present")
	assert.Equal(t, 44100.0, FixedPointToFloat32(uint32(sampleRates["alac"][0])), "Expected sample rate to be 44100.0 Hz")
}
