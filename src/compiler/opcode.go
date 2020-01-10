package compiler

type OpCode uint8

const (
	Constant OpCode = iota
	False
	True
	Nil

	Pop

	GetLocal
	SetLocal
	DefineGlobal
	GetGlobal
	SetGlobal
	GetUpvalue
	SetUpvalue
	GetProperty
	SetProperty
	GetSubscript
	SetSubscript

	Equal
	Greater
	GreaterEqual
	Less
	LessEqual
	NotEqual

	Not
	Negate

	Add
	Divide
	Exponentiate
	Multiply
	Reminder
	Subtract

	Return
)
