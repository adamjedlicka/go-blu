package value

type Nil struct{}

func (n Nil) IsTruthy() Boolean {
	return false
}

func (n Nil) ToString() String {
	return "nil"
}

func (n Nil) String() string {
	return "nil"
}
