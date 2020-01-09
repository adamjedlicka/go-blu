package compiler

type OpCode uint8

const (
	Constant OpCode = iota
	False
	True
	Nil

	Pop

	Equal
	Greater
	GreaterEqual
	Less
	LessEqual
	NotEqual

	Add
	Divide
	Exponentiate
	Multiply
	Reminder
	Subtract

	Return
)
