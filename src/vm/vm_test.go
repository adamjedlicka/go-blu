package vm

import (
	"github.com/adamjedlicka/go-blu/src/value"
	"testing"
)

func TestItHasBinaryNumberOperators(t *testing.T) {
	tests := []struct {
		code   string
		expect float64
	}{
		{
			code:   "1 + 1",
			expect: 2,
		},
		{
			code:   "1.1 + 1.1",
			expect: 2.2,
		},
		{
			code:   "100000000 + 1000",
			expect: 100001000,
		},
		{
			code:   "10 * 10",
			expect: 100,
		},
		{
			code:   "10 / 3",
			expect: float64(10) / float64(3),
		},
		{
			code:   "10 % 3",
			expect: 1,
		},
		{
			code:   "2^3",
			expect: 8,
		},
		{
			code:   "1 - 10",
			expect: -9,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := Exec(test.code)

			if float64(res.(value.Number)) != test.expect {
				t.Error(res)
			}
		})
	}
}

func TestItHasEqualityOperator(t *testing.T) {
	tests := []struct {
		code   string
		expect bool
	}{
		{
			code:   "1 == 1",
			expect: true,
		},
		{
			code:   "1.1 == 1.1",
			expect: true,
		},
		{
			code:   "true == true",
			expect: true,
		},
		{
			code:   "false == false",
			expect: true,
		},
		{
			code:   "1 == 2",
			expect: false,
		},
		{
			code:   "1.1 == 1.11",
			expect: false,
		},
		{
			code:   "true == false",
			expect: false,
		},
		{
			code:   "1 == true",
			expect: false,
		},
		{
			code:   "0 == false",
			expect: false,
		},
		{
			code:   "\"ab\" == \"a\" + \"b\"",
			expect: true,
		},
		{
			code:   "\"ab-1\" == \"a\" + \"b\" + - 1",
			expect: true,
		},
		{
			code:   "\"a\" == \"a\" + \"b\"",
			expect: false,
		},
		{
			code:   "!!true == !false",
			expect: true,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := Exec(test.code)

			if bool(res.(value.Boolean)) != test.expect {
				t.Error(res)
			}
		})
	}
}

func TestItHasInEqualityOperator(t *testing.T) {
	tests := []struct {
		code   string
		expect bool
	}{
		{
			code:   "return 1 != 1",
			expect: false,
		},
		{
			code:   "return false != true",
			expect: true,
		},
		{
			code:   "return false != nil",
			expect: true,
		},
		{
			code:   "return \"a\" != \"b\"",
			expect: true,
		},
		{
			code:   "return \"a\" != \"a\"",
			expect: false,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := Exec(test.code)

			if bool(res.(value.Boolean)) != test.expect {
				t.Error(res)
			}
		})
	}
}

func TestItHasComparisonOperators(t *testing.T) {
	tests := []struct {
		code   string
		expect bool
	}{
		{
			code:   "return 1 < 2",
			expect: true,
		},
		{
			code:   "return 1 < 1",
			expect: false,
		},
		{
			code:   "return 2 > 1",
			expect: true,
		},
		{
			code:   "return 2 >= 2",
			expect: true,
		},
		{
			code:   "return 3 <= 3",
			expect: true,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := Exec(test.code)

			if bool(res.(value.Boolean)) != test.expect {
				t.Error(res)
			}
		})
	}
}

func TestItHasGrouping(t *testing.T) {
	if Exec("(1 + 2) * 3").(value.Number) != 9 {
		t.Error("(1 + 2) * 3 != 9")
	}

	if Exec("1 + (2 * 3)").(value.Number) != 7 {
		t.Error("1 + (2 * 3) != 7")
	}

	if Exec("(1 + 2 * 3)").(value.Number) != 7 {
		t.Error("(1 + 2 * 3) != 7")
	}

	if Exec("1 + 2 * 3").(value.Number) != 7 {
		t.Error("1 + 2 * 3 != 7")
	}
}

func TestItHasGlobalVariables(t *testing.T) {
	res := Exec("var a = 3; a = a * a; a")

	if res.(value.Number) != 9 {
		t.Error(res)
	}
}

func TestItHasLocalVariables(t *testing.T) {
	if Exec("{ var a = 3; return a }").(value.Number) != 3 {
		t.Error("{ var a = 3; return a } should return 3")
	}

	if Exec("{ var a = 2; { var a = 3; a = 4 } return a; }").(value.Number) != 2 {
		t.Error("{ var a = 2; { var a = 3; a = 4 } return a; } should return 2")
	}
}

func TestItHasIfStatement(t *testing.T) {
	if Exec("if true { return 1; } return 2").(value.Number) != 1 {
		t.Error("not 1")
	}

	if Exec("if false { return 1; } return 2").(value.Number) != 2 {
		t.Error("not 2")
	}

	if Exec("if true: return 3; return 4").(value.Number) != 3 {
		t.Error("not 3")
	}

	if Exec("if false: return 3; return 4").(value.Number) != 4 {
		t.Error("not 4")
	}
}
