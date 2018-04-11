package internal

type BlockType uint16

func (bt BlockType) IsPlant() bool {
	if bt >= 17 && bt <= 31 {
		return true
	}
	return false
}

func (bt BlockType) IsTransparent() bool {
	if bt.IsPlant() {
		return true
	}
	switch bt {
	case 0, 10, 15:
		return true
	default:
		return false
	}
}

func (bt BlockType) IsObstacle() bool {
	if bt.IsPlant() {
		return false
	}
	switch bt {
	case 0:
		return false
	default:
		return true
	}
}
