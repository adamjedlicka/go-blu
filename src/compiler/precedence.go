package compiler

type Precedence uint8

type ParseFn func(c *Compiler, canAssign bool)

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

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

var parseRules = []ParseRule{
	{nil, nil, PrecedenceNone}, // At
	{nil, nil, PrecedenceNone}, // Caret
	{nil, nil, PrecedenceNone}, // Colon
	{nil, nil, PrecedenceNone}, // Comma
	{nil, nil, PrecedenceNone}, // Dot
	{nil, nil, PrecedenceNone}, // LeftBrace
	{nil, nil, PrecedenceNone}, // LeftBracket
	{nil, nil, PrecedenceNone}, // LeftParent
	{nil, nil, PrecedenceNone}, // Minus
	{nil, nil, PrecedenceNone}, // Percent
	{nil, nil, PrecedenceNone}, // Plus
	{nil, nil, PrecedenceNone}, // RightBrace
	{nil, nil, PrecedenceNone}, // RightBracket
	{nil, nil, PrecedenceNone}, // RightParen
	{nil, nil, PrecedenceNone}, // Semicolon
	{nil, nil, PrecedenceNone}, // Slash
	{nil, nil, PrecedenceNone}, // Start

	{nil, nil, PrecedenceNone}, // Bang
	{nil, nil, PrecedenceNone}, // BangEqual
	{nil, nil, PrecedenceNone}, // Equal
	{nil, nil, PrecedenceNone}, // EqualEqual
	{nil, nil, PrecedenceNone}, // Greater
	{nil, nil, PrecedenceNone}, // GreaterEqual
	{nil, nil, PrecedenceNone}, // Less
	{nil, nil, PrecedenceNone}, // LessEqual

	{nil, nil, PrecedenceNone}, // Identifier
	{nil, nil, PrecedenceNone}, // Number
	{nil, nil, PrecedenceNone}, // String

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
