package vm

import (
	"fmt"
	"github.com/adamjedlicka/go-blue/src/compiler"
	"github.com/adamjedlicka/go-blue/src/value"
	"math"
	"os"
)

type VM struct {
	chunk *compiler.Chunk
	stack *Stack
	ip    int

	globals map[value.String]value.Value
}

func NewVM() VM {
	return VM{
		chunk: nil,
		stack: nil,
		ip:    0,

		globals: make(map[value.String]value.Value),
	}
}

func Exec(source string) value.Value {
	vm := NewVM()

	return vm.Exec(source)
}

func (vm *VM) Exec(source string) value.Value {
	c := compiler.NewCompiler("script", source)
	chunk := c.Compile()
	if chunk == nil {
		return value.Nil{}
	}

	return vm.Interpret(chunk)
}

func (vm *VM) Interpret(chunk *compiler.Chunk) value.Value {
	vm.chunk = chunk
	vm.stack = NewStack()
	vm.ip = 0

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

		case compiler.GetLocal:
			panic("unimplemented")

		case compiler.SetLocal:
			panic("unimplemented")

		case compiler.DefineGlobal:
			name := vm.readString()

			vm.globals[name] = vm.Pop()

		case compiler.GetGlobal:
			name := vm.readString()

			if value, ok := vm.globals[name]; ok {
				vm.Push(value)
			} else {
				vm.runtimeError("Undefined global variable '%s'", name.ToString())
				return nil
			}

		case compiler.SetGlobal:
			name := vm.readString()

			if _, ok := vm.globals[name]; ok {
				vm.globals[name] = vm.Peek(0)
			} else {
				vm.runtimeError("Undefined global variable '%s'", name.ToString())
				return nil
			}

		case compiler.GetUpvalue:
			panic("unimplemented")

		case compiler.SetUpvalue:
			panic("unimplemented")

		case compiler.GetProperty:
			panic("unimplemented")

		case compiler.SetProperty:
			panic("unimplemented")

		case compiler.GetSubscript:
			panic("unimplemented")

		case compiler.SetSubscript:
			panic("unimplemented")

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

		case compiler.Not:
			vm.Push(!vm.Pop().IsTruthy())

		case compiler.Negate:
			vm.Push(-vm.Pop().(value.Number))

		case compiler.Add:
			rightValue := vm.Pop()
			leftValue := vm.Pop()

			switch left := leftValue.(type) {

			case value.Number:
				if right, ok := rightValue.(value.Number); ok {
					vm.Push(left + right)
				} else {
					vm.runtimeError("Both operands must be numbers.")
					return nil
				}

			case value.String:
				vm.Push(left + rightValue.ToString())

			default:
				vm.runtimeError("Left operand must be Number or String.")
				return nil
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
			panic("unreachable")
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

func (vm *VM) readConstant() value.Value {
	return vm.chunk.Constants()[vm.readShort()]
}

func (vm *VM) readString() value.String {
	return vm.readConstant().(value.String)
}

func (vm *VM) runtimeError(message string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, message, a...)
	_, _ = fmt.Fprintf(os.Stderr, "\n")

	line := vm.chunk.Lines()[vm.ip-1]
	name := vm.chunk.Name()

	_, _ = fmt.Fprintf(os.Stderr, "[line %d] in %s\n", line, name)
}
