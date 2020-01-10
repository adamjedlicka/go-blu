package value

type Boolean bool

func (b Boolean) IsTruthy() bool {
	return b != false
}
