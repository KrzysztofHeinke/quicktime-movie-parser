package atoms

type TkhdAtom struct {
	Version          uint8
	Flags            [3]byte
	CreationTime     uint32
	ModificationTime uint32
	TrackID          uint32
	Reserved         uint32
	Duration         uint32
	Reserved2        [8]byte
	Layer            uint16
	AlternateGroup   uint16
	Volume           uint16
	Reserved3        uint16
	Matrix           [36]byte
	Width            uint32
	Height           uint32
}
