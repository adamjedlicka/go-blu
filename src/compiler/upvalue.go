package compiler

type Upvalue struct {
	// The index of the local variable or upvalue being captured from the enclosing function.
	index uint16
	// Whether the captured variable is a local or upvalue in the enclosing function.
	isLocal bool
}
