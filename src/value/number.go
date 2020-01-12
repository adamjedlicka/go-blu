package value

func NumberVal(number float64) Value {
	return Value{
		value:  uintptr(number),
		object: nil,
	}
}

func IsNumber(value Value) bool {
	return value.value&qNaN != qNaN
}

func AsNumber(value Value) float64 {
	return float64(value.value)
}
