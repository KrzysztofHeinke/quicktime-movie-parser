package parser

import (
	"bytes"
	"testing"

	"github.com/KrzysztofHeinke/quicktime-movie-parser/pkg/models/atoms"
	"github.com/stretchr/testify/assert"
)

// Example data for testing
var tkhdData = []byte{
	0x00, 0x00, 0x00, 0x5C,
	0x74, 0x6B, 0x68, 0x64,
	0x00,
	0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02,
	0x00, 0x00, 0x00, 0x03,
	0x00, 0x00, 0x00, 0x04,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x01,
	0x00, 0x00,
	0x00, 0x02, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x02, 0x00, 0x00, // Width
	0x00, 0x03, 0x00, 0x00, // Height
}

// TestReadAtomHeader tests the ReadAtomHeader function.
func TestReadAtomHeader(t *testing.T) {
	reader := bytes.NewReader(tkhdData)
	header, err := ReadAtomHeader(reader)
	assert.NoError(t, err, "Expected no error reading atom header")
	assert.Equal(t, "tkhd", header.GetType(), "Expected atom type to be 'tkhd'")
	assert.Equal(t, uint32(92), header.GetSize(), "Expected atom size to be 92")
}

// TestFixedPointToFloat32 tests the FixedPointToFloat32 function.
func TestFixedPointToFloat32(t *testing.T) {
	width := uint32(0x00020000)  // 2.0 in Q16.16
	height := uint32(0x00030000) // 3.0 in Q16.16

	assert.Equal(t, 2.0, fixedPointToFloat32(width), "Expected width to be 2.0")
	assert.Equal(t, 3.0, fixedPointToFloat32(height), "Expected height to be 3.0")
}

// TestCreateTreeOfAtoms tests the CreateTreeOfAtoms function.
func TestCreateTreeOfAtoms(t *testing.T) {
	reader := bytes.NewReader(tkhdData)
	root, err := CreateTreeOfAtoms(reader)
	assert.NoError(t, err, "Expected no error creating tree of atoms")
	assert.NotNil(t, root, "Expected root to be not nil")
}

// TestCollectTrackInfo tests the CollectTrackInfo function.
func TestCollectTrackInfo(t *testing.T) {

	tkhdAtom := &atoms.LeafAtom{
		AtomHeader: atoms.AtomHeader{Size: 92, Type: [4]byte{'t', 'k', 'h', 'd'}},
		Data:       &atoms.TkhdAtom{Width: 0x00020000, Height: 0x00030000},
	}
	compositeAtom := &atoms.CompositeAtom{
		AtomHeader: atoms.AtomHeader{Type: [4]byte{'t', 'r', 'a', 'k'}},
	}
	compositeAtom.AddChild(tkhdAtom)
	root := &atoms.CompositeAtom{}
	root.AddChild(compositeAtom)

	CollectTrackInfo(root)
}

// TestCleanEmptyHeaders tests the CleanEmptyHeaders function.
func TestCleanEmptyHeaders(t *testing.T) {
	emptyCompositeAtom := &atoms.CompositeAtom{
		AtomHeader: atoms.AtomHeader{},
	}
	childAtom := &atoms.LeafAtom{
		AtomHeader: atoms.AtomHeader{Type: [4]byte{'c', 'h', 'l', 'd'}},
	}
	emptyCompositeAtom.AddChild(childAtom)
	root := &atoms.CompositeAtom{}
	root.AddChild(emptyCompositeAtom)

	cleanedRoot := CleanEmptyHeaders(root)
	assert.Equal(t, 1, len(cleanedRoot.(*atoms.CompositeAtom).GetChildren()), "Expected 1 child after cleaning")
	assert.Equal(t, childAtom, cleanedRoot.(*atoms.CompositeAtom).GetChildren()[0], "Expected the child to be preserved")
}
