package compiler

type Chunk struct {
	code []uint8
}

func NewChunk() Chunk {
	return Chunk{
		code: make([]uint8, 0),
	}
}

func (c *Chunk) pushOpCode(code OpCode) {
	c.code = append(c.code, uint8(code))
}
