package value

type String string

func (s String) IsTruthy() Boolean {
	return s != ""
}

func (s String) ToString() String {
	return s
}

func (s String) String() string {
	return string(s)
}
