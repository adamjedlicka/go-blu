package compiler

import "github.com/adamjedlicka/go-blu/src/parser"

type Precedence uint8

type ParseFn func(c *Compiler, canAssign bool)

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

type ParseRules = []ParseRule

const (
	PrecedenceNone       Precedence = iota
	PrecedenceAssignment            // =
	PrecedenceOr                    // or
	PrecedenceAnd                   // and
	PrecedenceEquality              // == !=
	PrecedenceComparison            // < > <= >=
	PrecedenceTerm                  // + -
	PrecedenceFactor                // * /
	PrecedencePower                 // ^
	PrecedenceUnary                 // ! -
	PrecedenceCall                  // . () []
	PrecedencePrimary
)

var parseRules ParseRules

func init() {
	parseRules = ParseRules{
		{nil, nil, PrecedenceNone},                              // At
		{nil, (*Compiler).binary, PrecedencePower},              // Caret
		{nil, nil, PrecedenceNone},                              // Colon
		{nil, nil, PrecedenceNone},                              // Comma
		{nil, nil, PrecedenceNone},                              // Dot
		{nil, nil, PrecedenceNone},                              // LeftBrace
		{nil, nil, PrecedenceNone},                              // LeftBracket
		{nil, nil, PrecedenceNone},                              // LeftParent
		{(*Compiler).unary, (*Compiler).binary, PrecedenceTerm}, // Minus
		{nil, (*Compiler).binary, PrecedenceFactor},             // Percent
		{nil, (*Compiler).binary, PrecedenceTerm},               // Plus
		{nil, nil, PrecedenceNone},                              // RightBrace
		{nil, nil, PrecedenceNone},                              // RightBracket
		{nil, nil, PrecedenceNone},                              // RightParen
		{nil, nil, PrecedenceNone},                              // Semicolon
		{nil, (*Compiler).binary, PrecedenceFactor},             // Slash
		{nil, (*Compiler).binary, PrecedenceFactor},             // Star

		{(*Compiler).unary, nil, PrecedenceNone},        // Bang
		{nil, (*Compiler).binary, PrecedenceEquality},   // BangEqual
		{nil, (*Compiler).binary, PrecedenceNone},       // Equal
		{nil, (*Compiler).binary, PrecedenceEquality},   // EqualEqual
		{nil, (*Compiler).binary, PrecedenceComparison}, // Greater
		{nil, (*Compiler).binary, PrecedenceComparison}, // GreaterEqual
		{nil, (*Compiler).binary, PrecedenceComparison}, // Less
		{nil, (*Compiler).binary, PrecedenceComparison}, // LessEqual

		{(*Compiler).variable, nil, PrecedenceNone}, // Identifier
		{(*Compiler).number, nil, PrecedenceNone},   // Number
		{(*Compiler).string, nil, PrecedenceNone},   // String

		{nil, nil, PrecedenceNone},                 // And
		{nil, nil, PrecedenceNone},                 // Assert
		{nil, nil, PrecedenceNone},                 // Break
		{nil, nil, PrecedenceNone},                 // Class
		{nil, nil, PrecedenceNone},                 // Echo
		{nil, nil, PrecedenceNone},                 // Else
		{(*Compiler).literal, nil, PrecedenceNone}, // False
		{nil, nil, PrecedenceNone},                 // Fn
		{nil, nil, PrecedenceNone},                 // For
		{nil, nil, PrecedenceNone},                 // Foreign
		{nil, nil, PrecedenceNone},                 // If
		{nil, nil, PrecedenceNone},                 // Import
		{(*Compiler).literal, nil, PrecedenceNone}, // Nil
		{nil, nil, PrecedenceNone},                 // Or
		{nil, nil, PrecedenceNone},                 // Return
		{(*Compiler).literal, nil, PrecedenceNone}, // True
		{nil, nil, PrecedenceNone},                 // Var
		{nil, nil, PrecedenceNone},                 // While

		{nil, nil, PrecedenceNone}, // Eof
		{nil, nil, PrecedenceNone}, // Newline

		{nil, nil, PrecedenceNone}, // Error
	}

	if len(parseRules)-1 != int(parser.Error) {
		panic("ParseRules table corrupt.")
	}
}
