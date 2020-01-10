package value

import (
	"strconv"
)

type Number float64

func (n Number) IsTruthy() Boolean {
	return n != 0
}

func (n Number) ToString() String {
	return String(strconv.FormatFloat(float64(n), 'f', -1, 64))
}
