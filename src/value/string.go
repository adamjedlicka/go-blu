package value

type String string

func StringVal(str string) Value {
	return ObjectVal(String(str))
}

func (s String) IsTruthy() bool {
	return s != ""
}

func (s String) ToString() string {
	return string(s)
}
