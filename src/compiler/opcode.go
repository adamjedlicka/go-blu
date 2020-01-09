package compiler

type OpCode uint8

const (
	Constant OpCode = iota
	False
	True
	Nil

	Pop

	Add

	Return
)
