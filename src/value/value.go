package value

import "unsafe"

type Value uintptr

// A mask that selects the sign bit.
const signBit = Value(1 << 63)

// The bits that must be set to indicate a quiet NaN.
const qNaN = Value(0x7ffc000000000000)

// Tag values for the different singleton values.
const tagNil = 1   // 01
const tagFalse = 2 // 10
const tagTrue = 3  // 11

func NilVal() Value {
	return qNaN | tagNil
}

func FalseVal() Value {
	return qNaN | tagFalse
}

func TrueVal() Value {
	return qNaN | tagTrue
}

func BooleanVal(boolean bool) Value {
	if boolean {
		return TrueVal()
	} else {
		return FalseVal()
	}
}

func NumberVal(number float64) Value {
	return Value(number)
}

func ObjectVal(object Object) Value {
	return signBit | qNaN | Value(unsafe.Pointer(&object))
}

func IsNil(value Value) bool {
	return value == NilVal()
}

func IsBool(value Value) bool {
	return (value & FalseVal()) == FalseVal()
}

func IsNumber(value Value) bool {
	return (value & qNaN) != qNaN
}

func IsObject(value Value) bool {
	return (value & (qNaN | signBit)) == (qNaN | signBit)
}

func AsBoolean(value Value) bool {
	return value == TrueVal()
}

func AsNumber(value Value) float64 {
	return float64(value)
}

func AsObject(value Value) Object {
	return *(*Object)(unsafe.Pointer(value))
}

func IsTruthy(value Value) bool {
	if IsNil(value) || AsBoolean(value) == false {
		return false
	}

	return true
}
