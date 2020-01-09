package vm

import (
	"github.com/adamjedlicka/go-blue/src/compiler"
	"github.com/adamjedlicka/go-blue/src/value"
)

type VM struct {
	chunk *compiler.Chunk
	stack *Stack
	ip    int
}

func NewVM() VM {
	return VM{
		chunk: nil,
		stack: NewStack(),
		ip:    0,
	}
}

func Exec(source string) value.Value {
	c := compiler.NewCompiler(source)
	chunk := c.Compile()
	vm := NewVM()
	return vm.Interpret(&chunk)
}

func (vm *VM) Interpret(chunk *compiler.Chunk) value.Value {
	vm.chunk = chunk

	for true {
		byte := vm.readByte()

		switch compiler.OpCode(byte) {

		case compiler.Constant:
			offset := vm.readShort()
			constant := vm.chunk.Constants()[offset]

			vm.stack.Push(constant)

		case compiler.False:
			vm.stack.Push(value.Boolean(false))

		case compiler.True:
			vm.stack.Push(value.Boolean(true))

		case compiler.Nil:
			vm.stack.Push(value.Nil{})

		case compiler.Pop:
			vm.stack.Pop()

		case compiler.Add:
			left := vm.stack.Pop().(value.Number)
			right := vm.stack.Pop().(value.Number)

			vm.stack.Push(left + right)

		case compiler.Return:
			return vm.stack.Pop()

		default:
			panic("Unknown OpCode")
		}
	}

	return nil
}

func (vm *VM) readByte() uint8 {
	byte := vm.chunk.Code()[vm.ip]

	vm.ip++

	return byte
}

func (vm *VM) readShort() uint16 {
	short1 := uint16(vm.chunk.Code()[vm.ip])
	short2 := uint16(vm.chunk.Code()[vm.ip+1])

	vm.ip += 2

	return (short1 << 8) | short2
}
