package atoms

type AtomIf interface {
	SetHeader(*AtomHeader)
	GetSize() uint32
	GetType() string
}

type AtomHeader struct {
	Size uint32
	Type [4]byte
}

func (a *AtomHeader) SetHeader(ah *AtomHeader) {
	*a = *ah
}

func (a *AtomHeader) GetType() string {
	return string(a.Type[:])
}

func (a *AtomHeader) GetSize() uint32 {
	return a.Size
}

type CompositeAtom struct {
	AtomHeader
	Childrens []AtomIf
}

func (ca *CompositeAtom) AddChild(atom AtomIf) {
	ca.Childrens = append(ca.Childrens, atom)
}

func (ca *CompositeAtom) GetChildren() []AtomIf {
	return ca.Childrens
}

func (ca *CompositeAtom) SetChildren(childrens []AtomIf) {
	ca.Childrens = childrens
}

type LeafAtom struct {
	AtomHeader
	Data any
}

func (la *LeafAtom) SetData(data any) {
	la.Data = data
}
