package vm

import "github.com/adamjedlicka/go-blue/src/value"

type Stack struct {
	values []value.Value
}

func NewStack() *Stack {
	s := new(Stack)
	s.values = make([]value.Value, 0)

	return s
}

func (s *Stack) Push(val value.Value) {
	s.values = append(s.values, val)
}

func (s *Stack) Pop() value.Value {
	val := s.values[len(s.values)-1]

	s.values = s.values[:len(s.values)-1]

	return val
}

func (s *Stack) Peek(distance int) value.Value {
	return s.values[len(s.values)-1-distance]
}
