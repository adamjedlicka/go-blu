package value

func BooleanVal(boolean bool) Value {
	if boolean {
		return TrueVal()
	} else {
		return FalseVal()
	}
}

func FalseVal() Value {
	return Value{
		value:  tagFalse,
		object: nil,
	}
}

func TrueVal() Value {
	return Value{
		value:  tagTrue,
		object: nil,
	}
}

func IsBoolean(value Value) bool {
	return value.value&tagFalse == tagFalse
}

func AsBoolean(value Value) bool {
	return value.value == tagTrue
}
