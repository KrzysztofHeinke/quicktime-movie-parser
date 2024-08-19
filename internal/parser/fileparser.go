package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

const chunkSize = 4096

// Parse is starting point to start parsing file.
func Parse(p string) {
	metadata, err := ReadFileMetadata(p)
	if err != nil {
		logrus.Fatal("Failed to read metadata of file")
	}
	tree, err := CreateTreeOfAtoms(bytes.NewReader(metadata))
	if err != nil {
		logrus.Fatal("Failed to create tree of atoms")
	}
	tree = CleanEmptyHeaders(tree)
	CollectTrackInfo(tree)
}

// FindAtomInFile is seeking for the specified atom in file
func FindAtomInFile(r io.Reader, search []byte) (int64, error) {
	var offset int64
	tailLen := len(search) - 1
	chunk := make([]byte, chunkSize+tailLen)
	n, err := r.Read(chunk[tailLen:])
	if err != nil {
		return -1, err
	}
	idx := bytes.Index(chunk[tailLen:n+tailLen], search)
	for {
		if idx >= 0 {
			return offset + int64(idx), nil
		}
		if err == io.EOF {
			return -1, nil
		} else if err != nil {
			return -1, err
		}
		copy(chunk, chunk[chunkSize:])
		offset += chunkSize
		n, err = r.Read(chunk[tailLen:])
		idx = bytes.Index(chunk[:n+tailLen], search)
	}
}

// ReadData taking a file, cursor posistion and reads data
func ReadData(f *os.File, cursorPosition int64, sizeToRead uint32) ([]byte, error) {
	var size []byte = make([]byte, sizeToRead)
	_, err := f.Seek(cursorPosition, io.SeekStart)
	if err != nil {
		return nil, err
	}
	_, err = f.Read(size[:sizeToRead])
	if err != nil {
		return nil, err
	}
	return size, nil
}

// ReadFileMetadata looks for the moov atom as the starting point and reads its contents.
func ReadFileMetadata(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		logrus.Errorf("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	position, err := FindAtomInFile(file, []byte("moov"))
	if err != nil {
		logrus.Errorf("Error finding moov atom:", err)
		return nil, err
	}

	// Go to the start of the "moov" atom minus header and moov size
	_, err = file.Seek(position-7, io.SeekStart)
	if err != nil {
		logrus.Errorf("Error seeking to moov atom:", err)
		return nil, err
	}

	moovSizeData := make([]byte, 4)
	_, err = file.Read(moovSizeData)
	if err != nil {
		logrus.Errorf("Error reading moov size:", err)
		return nil, err
	}
	size := binary.BigEndian.Uint32(moovSizeData)

	// Read the entire "moov" atom (including its header)
	atomData, err := ReadData(file, position-7, size)
	if err != nil {
		logrus.Errorf("Error reading moov atom data:", err)
		return nil, err
	}

	return atomData, nil
}

// SearchAtoms is checking if in file there is specific atom
func SearchAtoms(data []byte, searchBytes []byte) (int, error) {
	index := bytes.Index(data, searchBytes)

	if index != -1 {
		fmt.Printf("Found '%s' at position: %d\n", searchBytes, index)
	} else {
		fmt.Printf("'%s' not found", searchBytes)
		return -1, fmt.Errorf("%s not found", searchBytes)
	}
	return index, nil
}
