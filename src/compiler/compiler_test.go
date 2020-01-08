package compiler

import "testing"

func TestItCompilesLiterals(t *testing.T) {
	compiler := NewCompiler("   true   ;  false   ")
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

	if chunk.code[3] != uint8(Pop) {
		t.Error(chunk.code[3])
	}
}
