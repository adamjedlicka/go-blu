package value

import "strconv"

type Value struct {
	value  uintptr
	object Object
}

// A mask that selects the sign bit.
const signBit = 1 << 63

// The bits that must be set to indicate a quiet NaN.
const qNaN = 0x7ffc000000000000

// Tag values for the different singleton values.
const tagNil = qNaN | 1
const tagFalse = qNaN | 2
const tagTrue = qNaN | 3

func IsTruthy(value Value) bool {
	if IsNil(value) || AsBoolean(value) == false {
		return false
	}

	return true
}

func (v Value) String() string {
	if IsNil(v) {
		return "nil"
	} else if IsBoolean(v) {
		if AsBoolean(v) {
			return "true"
		} else {
			return "false"
		}
	} else if IsNumber(v) {
		return strconv.FormatFloat(float64(v.value), 'f', -1, 64)
	} else {
		return v.object.ToString()
	}
}
