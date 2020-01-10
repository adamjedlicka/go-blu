package value

type Number float64

func (n Number) IsTruthy() bool {
	return n != 0
}
