package value

type String string

func (s String) IsTruthy() Boolean {
	return s != ""
}

func StringVal(str string) Value {
	return ObjectVal(String(str))
}

func (s String) ToString() String {
	return s
}

func (s String) String() string {
	return string(s)
}
