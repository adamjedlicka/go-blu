package compiler

import "github.com/adamjedlicka/go-blu/src/value"

const MaxConstants = 65000

type Chunk struct {
	name string

	code      []uint8
	lines     []int
	constants []value.Value
}

func NewChunk(name string) *Chunk {
	return &Chunk{
		name: name,

		code:      make([]uint8, 0),
		lines:     make([]int, 0),
		constants: make([]value.Value, 0),
	}
}

func (c *Chunk) Name() string {
	return c.name
}

func (c *Chunk) Code() []uint8 {
	return c.code
}

func (c *Chunk) Lines() []int {
	return c.lines
}

func (c *Chunk) Constants() []value.Value {
	return c.constants
}

func (c *Chunk) pushCode(code uint8, line int) {
	c.code = append(c.code, code)
	c.lines = append(c.lines, line)
}

func (c *Chunk) pushConstant(constant value.Value) uint16 {
	if len(c.constants) == MaxConstants {
		panic("Too many constants in one chunk.")
	}

	c.constants = append(c.constants, constant)

	return uint16(len(c.constants) - 1)
}
