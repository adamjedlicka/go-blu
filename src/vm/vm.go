package vm

import (
	"github.com/adamjedlicka/go-blue/src/compiler"
	"github.com/adamjedlicka/go-blue/src/value"
	"math"
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
		switch compiler.OpCode(vm.readByte()) {

		case compiler.Constant:
			offset := vm.readShort()
			constant := vm.chunk.Constants()[offset]

			vm.Push(constant)

		case compiler.False:
			vm.Push(value.Boolean(false))

		case compiler.True:
			vm.Push(value.Boolean(true))

		case compiler.Nil:
			vm.Push(value.Nil{})

		case compiler.Pop:
			vm.Pop()

		case compiler.Equal:
			left := vm.Pop()
			right := vm.Pop()

			vm.Push(value.Boolean(left == right))

		case compiler.Greater:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Boolean(left > right))

		case compiler.GreaterEqual:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Boolean(left >= right))

		case compiler.Less:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Boolean(left < right))

		case compiler.LessEqual:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Boolean(left <= right))

		case compiler.NotEqual:
			right := vm.Pop()
			left := vm.Pop()

			vm.Push(value.Boolean(left != right))

		case compiler.Add:
			rightValue := vm.Pop()
			leftValue := vm.Pop()

			switch left := leftValue.(type) {

			case value.Number:
				if right, ok := rightValue.(value.Number); ok {
					vm.Push(left + right)
				} else {
					// TODO : Error handling
					panic("Both operands must be numbers.")
				}

			case value.String:
				vm.Push(left + rightValue.ToString())

			default:
				// TODO : Error handling
				panic("Left operand must be Number or String.")
			}

		case compiler.Divide:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(left / right)

		case compiler.Exponentiate:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Number(math.Pow(float64(left), float64(right))))

		case compiler.Multiply:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(left * right)

		case compiler.Reminder:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(value.Number(int(left) % int(right)))

		case compiler.Subtract:
			right := vm.Pop().(value.Number)
			left := vm.Pop().(value.Number)

			vm.Push(left - right)

		case compiler.Return:
			return vm.Pop()

		default:
			panic("Unknown OpCode")
		}
	}

	return nil
}

func (vm *VM) Push(val value.Value) {
	vm.stack.Push(val)
}

func (vm *VM) Pop() value.Value {
	return vm.stack.Pop()
}

func (vm *VM) Peek(distance int) value.Value {
	return vm.stack.Peek(distance)
}

func (vm *VM) readByte() uint8 {
	vm.ip++

	return vm.chunk.Code()[vm.ip-1]
}

func (vm *VM) readShort() uint16 {
	short1 := uint16(vm.chunk.Code()[vm.ip])
	short2 := uint16(vm.chunk.Code()[vm.ip+1])

	vm.ip += 2

	return (short1 << 8) | short2
}
