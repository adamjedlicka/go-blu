package value

type Value interface {
	IsTruthy() Boolean
	ToString() String
}
