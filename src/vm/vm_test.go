package vm

import (
	"github.com/adamjedlicka/go-blue/src/value"
	"testing"
)

func TestItAddsNumbers(t *testing.T) {
	res := Exec("return 1+1")

	if res.(value.Number) != 2 {
		t.Error(res)
	}
}
