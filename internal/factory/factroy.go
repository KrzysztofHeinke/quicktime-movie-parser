package factory

import (
	"bytes"
	"encoding/binary"

	"github.com/KrzysztofHeinke/quicktime-movie-parser/pkg/models/atoms"
	"github.com/sirupsen/logrus"
)

func AtomFactory(header atoms.AtomHeader, reader *bytes.Reader) (any, error) {
	switch header.GetType() {
	case "moov":
	case "mvhd":
	case "clip":
	case "crgn":
	case "udta":
	case "tkhd":
		result := &atoms.TkhdAtom{}
		CastToStruct(reader, result)
		return result, nil
	case "matt":
	case "kmat":
	case "edts":
	case "trak":
	case "mdhd":
		result := &atoms.MdhdAtom{}
		CastToStruct(reader, result)
		return result, nil
	case "hdlr":
	case "minf":
	case "elst":
	case "vmhd":
	case "dref":
	case "stts":
	case "stss":
	case "stsd":
		return atoms.ParseStsdAtom(reader)
	case "stsz":
	case "stsc":
	case "stco":
	default:

	}
	return nil, nil
}

func CastToStruct(reader *bytes.Reader, structToCast any) any {
	err := binary.Read(reader, binary.BigEndian, structToCast)
	if err != nil {
		logrus.Fatalf("Failed to read binary data into struct: %v", err)
	}
	return structToCast
}
