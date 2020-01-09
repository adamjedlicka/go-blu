package value

type Value interface {
	IsTruthy() bool
}

type Number float64

func (n Number) IsTruthy() bool {
	return n != 0
}

type Boolean bool

func (b Boolean) IsTruthy() bool {
	return b != false
}

type Nil struct{}

func (n Nil) IsTruthy() bool {
	return false
}
