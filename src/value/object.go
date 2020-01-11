package value

type Object interface {
	IsTruthy() Boolean
	ToString() String
}
