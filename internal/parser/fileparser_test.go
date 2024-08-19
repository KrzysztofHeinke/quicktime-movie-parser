package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFindAtomInFile tests the FindAtomInFile function
func TestFindAtomInFile(t *testing.T) {
	testData := []byte("randomdataandmoovatomdata")

	reader := bytes.NewReader(testData)
	position, err := FindAtomInFile(reader, []byte("moov"))
	assert.NoError(t, err, "Expected no error while searching for 'moov'")
	assert.Equal(t, int64(13), position, "Expected 'moov' atom to be found at position 13")

	reader = bytes.NewReader(testData)
	position, err = FindAtomInFile(reader, []byte("notfound"))
	assert.NoError(t, err, "Expected no error while searching for a non-existent atom")
	assert.Equal(t, int64(-1), position, "Expected -1 position for a non-existent atom")
}

// TestReadData tests the ReadData function
func TestReadData(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("This is some test data for reading."))
	assert.NoError(t, err, "Expected no error writing to temp file")
	_, err = tmpFile.Seek(0, 0)
	assert.NoError(t, err, "Expected no error while setting the cursor")

	data, err := ReadData(tmpFile, 5, 9)
	assert.NoError(t, err, "Expected no error reading data from file")
	assert.Equal(t, []byte("is some t"), data, "Expected to read 'is some t' from file")
}

// TestSearchAtoms tests the SearchAtoms function
func TestSearchAtoms(t *testing.T) {
	data := []byte("randomdataandmoovatomdata")

	position, err := SearchAtoms(data, []byte("moov"))
	assert.NoError(t, err, "Expected no error while searching for 'moov'")
	assert.Equal(t, 13, position, "Expected 'moov' atom to be found at position 13")

	position, err = SearchAtoms(data, []byte("notfound"))
	assert.Error(t, err, "Expected error while searching for non-existent atom")
	assert.Equal(t, -1, position, "Expected position -1 for non-existent atom")
}
