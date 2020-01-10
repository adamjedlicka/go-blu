package value

type Nil struct{}

func (n Nil) IsTruthy() bool {
	return false
}
