package value

type Boolean bool

func (b Boolean) IsTruthy() Boolean {
	return b != false
}

func (b Boolean) ToString() String {
	if b {
		return "true"
	}

	return "false"
}
