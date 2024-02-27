package DomainEntity

type Box struct {
	boxType     string
	maxQuantity int
}

func NewBox(
	boxType string,
	maxQuantity int,
) *Box {
	box := new(Box)
	box.boxType = boxType
	box.maxQuantity = maxQuantity
	return box
}

func (box *Box) GetBoxType() string {
	return box.boxType
}

func (box *Box) GetMaxQuantity() int {
	return box.maxQuantity
}
