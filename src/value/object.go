package value

type Object interface {
	IsTruthy() bool
	ToString() string
}

func ObjectVal(object Object) Value {
	return Value{
		value:  qNaN | signBit,
		object: object,
	}
}

func IsObject(value Value) bool {
	return value.object != nil
}

func AsObject(value Value) Object {
	return value.object
}
