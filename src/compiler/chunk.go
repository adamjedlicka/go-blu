package compiler

import "github.com/adamjedlicka/go-blue/src/value"

const MaxConstants = 65000

type Chunk struct {
	code      []uint8
	constants []value.Value
}

func NewChunk() *Chunk {
	return &Chunk{
		code:      make([]uint8, 0),
		constants: make([]value.Value, 0),
	}
}

func (c *Chunk) Code() []uint8 {
	return c.code
}

func (c *Chunk) Constants() []value.Value {
	return c.constants
}

func (c *Chunk) pushCode(code uint8) {
	c.code = append(c.code, code)
}

func (c *Chunk) pushConstant(constant value.Value) uint16 {
	if len(c.constants) == MaxConstants {
		panic("Too many constants in one chunk.")
	}

	c.constants = append(c.constants, constant)

	return uint16(len(c.constants) - 1)
}
