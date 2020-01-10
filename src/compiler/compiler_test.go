package compiler

import (
	"github.com/adamjedlicka/go-blu/src/value"
	"testing"
)

func TestItCompilesLiterals(t *testing.T) {
	compiler := NewCompiler("", "   true   ;  false   ")
	chunk := compiler.Compile()

	if chunk.code[0] != uint8(True) {
		t.Error(chunk.code[0])
	}

	if chunk.code[1] != uint8(Pop) {
		t.Error(chunk.code[1])
	}

	if chunk.code[2] != uint8(False) {
		t.Error(chunk.code[2])
	}

	if chunk.code[3] != uint8(Return) {
		t.Error(chunk.code[3])
	}
}

func TestItCompilesNumbers(t *testing.T) {
	compiler := NewCompiler("", "123.456; 789")
	chunk := compiler.Compile()

	if chunk.code[0] != uint8(Constant) {
		t.Error(chunk.code[0])
	}

	if chunk.code[1] != 0 {
		t.Error(chunk.code[1])
	}

	if chunk.code[2] != 0 {
		t.Error(chunk.code[2])
	}

	if chunk.constants[0].(value.Number) != 123.456 {
		t.Error(chunk.constants[0])
	}

	if chunk.code[3] != uint8(Pop) {
		t.Error(chunk.code[3])
	}

	if chunk.code[4] != uint8(Constant) {
		t.Error(chunk.code[4])
	}

	if chunk.code[5] != 0 {
		t.Error(chunk.code[5])
	}

	if chunk.code[6] != 1 {
		t.Error(chunk.code[6])
	}

	if chunk.constants[1].(value.Number) != 789 {
		t.Error(chunk.constants[1])
	}
}

func TestItCompilesBinaryOperators(t *testing.T) {
	compiler := NewCompiler("", "1 + 2")
	chunk := compiler.Compile()

	if chunk.code[6] != uint8(Add) {
		t.Error(chunk.code[6])
	}

	if chunk.constants[0].(value.Number) != 1 {
		t.Error(chunk.constants[0])
	}

	if chunk.constants[1].(value.Number) != 2 {
		t.Error(chunk.constants[1])
	}
}
