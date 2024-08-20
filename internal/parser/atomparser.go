package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/KrzysztofHeinke/quicktime-movie-parser/internal/factory"
	"github.com/KrzysztofHeinke/quicktime-movie-parser/pkg/models/atoms"
	"github.com/sirupsen/logrus"
)

// CreateTreeOfAtoms parses atoms from the reader and constructs a tree of atoms.
func CreateTreeOfAtoms(reader *bytes.Reader) (atoms.AtomIf, error) {
	root := &atoms.CompositeAtom{}
	dataSize := int64(reader.Len())
	dataRead := int64(0)

	for dataRead < dataSize-4 {
		startPos, err := reader.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("error getting start position: %w", err)
		}

		header, err := ReadAtomHeader(reader)
		if err != nil {
			return nil, err
		}
		if header == nil {
			break
		}

		atomType := header.GetType()
		atomSize := int64(header.GetSize())
		headerSize := int64(binary.Size(*header))

		if atomSize <= headerSize {
			return nil, fmt.Errorf("invalid atom size: %d (header size: %d)", atomSize, headerSize)
		}

		if isCompositeAtom(atomType) {
			logrus.Debugf("Found composite atom: %s", atomType)

			compositeAtom := &atoms.CompositeAtom{
				AtomHeader: *header,
			}

			remainingSize := atomSize - headerSize
			if remainingSize <= 0 {
				return nil, fmt.Errorf("invalid remaining size: %d", remainingSize)
			}

			sectionReader := io.NewSectionReader(reader, startPos+headerSize, remainingSize)
			sectionData, err := ReadBytes(sectionReader, int(remainingSize))
			if err != nil {
				return nil, err
			}

			childRoot, err := CreateTreeOfAtoms(bytes.NewReader(sectionData))
			if err != nil {
				return nil, err
			}
			compositeAtom.AddChild(childRoot)

			root.AddChild(compositeAtom)

			dataRead = startPos + atomSize
			if _, err := reader.Seek(dataRead, io.SeekStart); err != nil {
				return nil, fmt.Errorf("error seeking after composite atom: %w", err)
			}
		} else {
			logrus.Debugf("Found leaf atom: %s", atomType)
			atomData, err := ReadBytes(reader, int(atomSize))
			if err != nil {
				return nil, err
			}

			atomAdditionalData, err := factory.AtomFactory(*header, bytes.NewReader(atomData[8:]))
			if err != nil {
				return nil, err
			}

			leafAtom := &atoms.LeafAtom{
				AtomHeader: *header,
				Data:       atomAdditionalData,
			}
			root.AddChild(leafAtom)
			dataRead = startPos + atomSize
			if _, err := reader.Seek(dataRead, io.SeekStart); err != nil {
				return nil, fmt.Errorf("error seeking after leaf atom: %w", err)
			}
		}
	}

	return root, nil
}

// ReadAtomHeader reads the atom header from the reader.
func ReadAtomHeader(reader *bytes.Reader) (*atoms.AtomHeader, error) {
	header := atoms.AtomHeader{}
	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, fmt.Errorf("error during header reading: %w", err)
	}
	if _, err := reader.Seek(-int64(binary.Size(header)), io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("error seeking after reading header: %w", err)
	}

	return &header, nil
}

// isCompositeAtom returns true if the atom is a composite atom.
func isCompositeAtom(atomType string) bool {
	compositeAtoms := map[string]bool{
		"moov": true,
		"trak": true,
		"mdia": true,
		"minf": true,
		"stbl": true,
	}
	return compositeAtoms[atomType]
}

// ReadBytes reads the specified number of bytes from the reader.
func ReadBytes(reader io.Reader, size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size: %d", size)
	}

	buf := make([]byte, size)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read bytes: %v", err)
	}

	if n != size {
		return nil, fmt.Errorf("incomplete read: %d bytes read, %d bytes expected", n, size)
	}

	return buf, nil
}

// CleanEmptyHeaders recursively removes composite atoms with empty headers and moves their children up the tree.
func CleanEmptyHeaders(root atoms.AtomIf) atoms.AtomIf {
	if compositeAtom, ok := root.(*atoms.CompositeAtom); ok {
		var newChildren []atoms.AtomIf

		for _, child := range compositeAtom.GetChildren() {
			cleanedChild := CleanEmptyHeaders(child)

			if compositeChild, ok := cleanedChild.(*atoms.CompositeAtom); ok {
				if isHeaderEmpty(compositeChild.AtomHeader) && len(compositeChild.GetChildren()) > 0 {
					newChildren = append(newChildren, compositeChild.GetChildren()...)
				} else {
					newChildren = append(newChildren, cleanedChild)
				}
			} else {
				newChildren = append(newChildren, cleanedChild)
			}
		}
		compositeAtom.SetChildren(newChildren)
		if isHeaderEmpty(compositeAtom.AtomHeader) && len(compositeAtom.GetChildren()) > 0 {
			return &atoms.CompositeAtom{
				AtomHeader: compositeAtom.AtomHeader,
				Childrens:  compositeAtom.GetChildren(),
			}
		}
	}

	return root
}

// isHeaderEmpty checks if the given AtomHeader is considered empty.
func isHeaderEmpty(header atoms.AtomHeader) bool {
	return header.GetSize() == 0 || header.GetType() == ""
}

// CollectTrackInfo collects and prints track information from the atom tree.
func CollectTrackInfo(root atoms.AtomIf) {
	if compositeAtom, ok := root.(*atoms.CompositeAtom); ok {
		for _, child := range compositeAtom.GetChildren() {
			if trakAtom, ok := child.(*atoms.CompositeAtom); ok {
				var isVideoTrack bool
				var width, height float64

				for _, trakChild := range trakAtom.GetChildren() {
					switch atom := trakChild.(type) {
					case *atoms.CompositeAtom:
						CollectTrackInfo(atom)
					case *atoms.LeafAtom:
						if atom.AtomHeader.GetType() == "tkhd" {
							a := atom.Data.(*atoms.TkhdAtom)
							width = fixedPointToFloat32(a.Width)
							height = fixedPointToFloat32(a.Height)
							isVideoTrack = true
						}
						if atom.AtomHeader.GetType() == "stsd" {
							a := atom.Data.(*atoms.AtomStsd)
							sampleRates, _ := atoms.GetSampleRates(a)
							for codec, rates := range sampleRates {
								for _, rate := range rates {
									logrus.Infof("Codec: %s, Sample Rate: %.2f Hz\n", codec, fixedPointToFloat32(uint32(rate)))
								}
							}
						}
					}
				}
				if isVideoTrack {
					logrus.Infof("Video Track: Width = %.2f, Height = %.2f\n", width, height)
				}
			}
			CollectTrackInfo(child)
		}
	}
}

// fixedPointToFloat32 converts a fixed-point Q16.16 value to a floating-point number
func fixedPointToFloat32(value uint32) float64 {
	return float64(value) / (1 << 16)
}
