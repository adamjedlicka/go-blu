package value

import "testing"

func TestNilEqualsToItself(t *testing.T) {
	a := Nil{}
	b := Nil{}

	if a != b {
		t.Error("Nil does not equal to itself")
	}
}
