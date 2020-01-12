package vm

import (
	"fmt"
	"github.com/adamjedlicka/go-blu/src/compiler"
	"github.com/adamjedlicka/go-blu/src/parser"
	"github.com/adamjedlicka/go-blu/src/value"
	"math"
	"os"
)

const StackMax = 256

type VM struct {
	chunk *compiler.Chunk

	ip int

	stack    [StackMax]value.Value
	stackLen int

	globals map[value.String]value.Value
}

func NewVM() VM {
	return VM{
		chunk: nil,

		stack: [StackMax]value.Value{},
		ip:    0,

		globals: make(map[value.String]value.Value),
	}
}

func Exec(source string) value.Value {
	vm := NewVM()

	return vm.Exec(source)
}

func (vm *VM) Exec(source string) value.Value {
	p := parser.NewParser([]rune(source))
	c := compiler.NewCompiler("script", p)
	chunk := c.Compile()
	if chunk == nil {
		return value.NilVal()
	}

	return vm.Interpret(chunk)
}

func (vm *VM) Interpret(chunk *compiler.Chunk) value.Value {
	vm.chunk = chunk
	vm.ip = 0

	vm.stackLen = 0

	for true {
		switch compiler.OpCode(vm.readByte()) {

		case compiler.Constant:
			offset := vm.readShort()
			constant := vm.chunk.Constants()[offset]

			vm.Push(constant)

		case compiler.False:
			vm.Push(value.FalseVal())

		case compiler.True:
			vm.Push(value.TrueVal())

		case compiler.Nil:
			vm.Push(value.NilVal())

		case compiler.Pop:
			vm.Pop()

		case compiler.GetLocal:
			slot := vm.readShort()

			vm.Push(vm.stack[slot])

		case compiler.SetLocal:
			slot := vm.readShort()

			vm.stack[slot] = vm.Peek(0)

		case compiler.DefineGlobal:
			name := vm.readString()

			vm.globals[name] = vm.Pop()

		case compiler.GetGlobal:
			name := vm.readString()

			if val, ok := vm.globals[name]; ok {
				vm.Push(val)
			} else {
				vm.runtimeError("Undefined global variable '%s'", name.ToString())
				return value.NilVal()
			}

		case compiler.SetGlobal:
			name := vm.readString()

			if _, ok := vm.globals[name]; ok {
				vm.globals[name] = vm.Peek(0)
			} else {
				vm.runtimeError("Undefined global variable '%s'", name.ToString())
				return value.NilVal()
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

			vm.Push(value.BooleanVal(left == right))

		case compiler.Greater:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.BooleanVal(left > right))

		case compiler.GreaterEqual:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.BooleanVal(left >= right))

		case compiler.Less:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.BooleanVal(left < right))

		case compiler.LessEqual:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.BooleanVal(left <= right))

		case compiler.NotEqual:
			right := vm.Pop()
			left := vm.Pop()

			vm.Push(value.BooleanVal(left != right))

		case compiler.Not:
			vm.Push(value.BooleanVal(!value.IsTruthy(vm.Pop())))

		case compiler.Negate:
			vm.Push(value.NumberVal(-value.AsNumber(vm.Pop())))

		case compiler.Add:
			right := vm.Pop()
			left := vm.Pop()

			if value.IsNumber(left) && value.IsNumber(right) {
				vm.Push(value.NumberVal(value.AsNumber(left) + value.AsNumber(right)))
			} else if value.IsObject(left) && value.IsObject(right) {
				left := value.AsObject(left).(value.String)
				right := value.AsObject(right).(value.String)

				vm.Push(value.StringVal(string(left + right)))
			} else {
				vm.runtimeError("Both operands must be numbers.")
				return value.NilVal()
			}

		case compiler.Divide:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.NumberVal(left / right))

		case compiler.Exponentiate:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.NumberVal(math.Pow(left, right)))

		case compiler.Multiply:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.NumberVal(left * right))

		case compiler.Reminder:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.NumberVal(float64(int(left) % int(right))))

		case compiler.Subtract:
			right := value.AsNumber(vm.Pop())
			left := value.AsNumber(vm.Pop())

			vm.Push(value.NumberVal(left - right))

		case compiler.Jump:
			offset := vm.readShort()

			vm.ip += int(offset)

		case compiler.JumpIfFalsy:
			offset := vm.readShort()

			if !value.IsTruthy(vm.Peek(0)) {
				vm.ip += int(offset)
			}

		case compiler.JumpIfTruthy:
			offset := vm.readShort()

			if value.IsTruthy(vm.Peek(0)) {
				vm.ip += int(offset)
			}

		case compiler.Loop:
			offset := vm.readShort()

			vm.ip -= int(offset)

		case compiler.Return:
			return vm.Pop()

		default:
			panic("unreachable")
		}
	}

	return value.NilVal()
}

func (vm *VM) Push(val value.Value) {
	vm.stack[vm.stackLen] = val

	vm.stackLen++
}

func (vm *VM) Pop() value.Value {
	vm.stackLen--

	return vm.stack[vm.stackLen]
}

func (vm *VM) Peek(distance int) value.Value {
	return vm.stack[vm.stackLen-1-distance]
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
	return value.AsObject(vm.readConstant()).(value.String)
}

func (vm *VM) runtimeError(message string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, message, a...)
	_, _ = fmt.Fprintf(os.Stderr, "\n")

	line := vm.chunk.Lines()[vm.ip-1]
	name := vm.chunk.Name()

	_, _ = fmt.Fprintf(os.Stderr, "[line %d] in %s\n", line, name)
}
