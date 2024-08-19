package atoms

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// AtomStsd represents the 'stsd' atom structure in a media file
type AtomStsd struct {
	Version       byte
	Flags         [3]byte
	EntryCount    uint32
	SampleEntries []SampleEntry
}

// SampleEntry represents a sample entry in the 'stsd' atom
type SampleEntry struct {
	Size     [4]byte
	Type     [4]byte
	Reserved [6]byte
	RefIndex uint16
	Data     []byte
}

// ParseStsdAtom parses the 'stsd' atom and extracts details for audio streams
func ParseStsdAtom(reader io.Reader) (*AtomStsd, error) {
	var stsd AtomStsd

	if err := binary.Read(reader, binary.BigEndian, &stsd.Version); err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}
	if _, err := reader.Read(stsd.Flags[:]); err != nil {
		return nil, fmt.Errorf("error reading flags: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &stsd.EntryCount); err != nil {
		return nil, fmt.Errorf("error reading entry count: %w", err)
	}

	stsd.SampleEntries = make([]SampleEntry, stsd.EntryCount)
	for i := uint32(0); i < stsd.EntryCount; i++ {
		var entry SampleEntry
		if err := binary.Read(reader, binary.BigEndian, &entry.Size); err != nil {
			return nil, fmt.Errorf("error reading data size of data: %w", err)
		}
		if _, err := reader.Read(entry.Type[:]); err != nil {
			return nil, fmt.Errorf("error reading sample entry type: %w", err)
		}
		if _, err := reader.Read(entry.Reserved[:]); err != nil {
			return nil, fmt.Errorf("error reading reserved bytes: %w", err)
		}
		if err := binary.Read(reader, binary.BigEndian, &entry.RefIndex); err != nil {
			return nil, fmt.Errorf("error reading data reference index: %w", err)
		}

		entryHeaderSize := binary.Size(entry.Type) + binary.Size(entry.Reserved) + binary.Size(entry.RefIndex)
		entrySize := 78 // Default size for simplicity; adjust based on actual data structure
		entry.Data = make([]byte, entrySize-entryHeaderSize)
		if _, err := reader.Read(entry.Data); err != nil {
			return nil, fmt.Errorf("error reading sample entry data: %w", err)
		}

		stsd.SampleEntries[i] = entry
	}

	return &stsd, nil
}

// GetSampleRates extracts the sample rates for all audio sample entries
func GetSampleRates(stsd *AtomStsd) (map[string][]float64, error) {
	sampleRates := make(map[string][]float64)

	for _, entry := range stsd.SampleEntries {
		switch string(entry.Type[:]) {
		case "mp4a":
			if len(entry.Data) >= 16 {
				rate := binary.BigEndian.Uint32(entry.Data[16:20])
				sampleRates["mp4a"] = append(sampleRates["mp4a"], float64(rate))
			}
		case "ac-3", "ec-3", "alac":
			if len(entry.Data) >= 16 {
				rate := binary.BigEndian.Uint32(entry.Data[16:20])
				sampleRates[string(entry.Type[:])] = append(sampleRates[string(entry.Type[:])], float64(rate))
			}
		default:
			logrus.Debugf("Unsupported audio type: %s\n", string(entry.Type[:]))
		}
	}

	if len(sampleRates) == 0 {
		return nil, fmt.Errorf("no sample rates found for audio entries")
	}

	return sampleRates, nil
}
