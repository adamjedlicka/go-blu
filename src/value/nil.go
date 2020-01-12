package value

func NilVal() Value {
	return Value{
		value:  tagNil,
		object: nil,
	}
}

func IsNil(value Value) bool {
	return value.value == tagNil
}
