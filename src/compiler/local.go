package compiler

import "github.com/adamjedlicka/go-blu/src/parser"

type Local struct {
	name  parser.Token
	depth int8
	// True if this local variable is captured as an upvalue by a function.
	isUpvalue bool
}
