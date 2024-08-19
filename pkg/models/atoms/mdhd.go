package atoms

type MdhdAtom struct {
	Version          uint8
	Flags            [3]byte
	CreationTime     uint32
	ModificationTime uint32
	TimeScale        uint32
	Duration         uint32
	Language         uint16
	Quality          uint16
}
