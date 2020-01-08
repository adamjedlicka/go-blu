package parser

type TokenType string

const (
	// Single-character tokens
	At           TokenType = "At"
	Caret        TokenType = "Caret"
	Colon        TokenType = "Colon"
	Comma        TokenType = "Comma"
	Dot          TokenType = "Dot"
	LeftBrace    TokenType = "LeftBrace"
	LeftBracket  TokenType = "LeftBracket"
	LeftParen    TokenType = "LeftParen"
	Minus        TokenType = "Minus"
	Percent      TokenType = "Percent"
	Plus         TokenType = "Plus"
	RightBrace   TokenType = "RightBrace"
	RightBracket TokenType = "RightBracket"
	RightParen   TokenType = "RightParen"
	Semicolon    TokenType = "Semicolon"
	Slash        TokenType = "Slash"
	Star         TokenType = "Star"

	// One or two character tokens
	Bang         TokenType = "Bang"
	BangEqual    TokenType = "BangEqual"
	Equal        TokenType = "Equal"
	EqualEqual   TokenType = "EqualEqual"
	Greater      TokenType = "Greater"
	GreaterEqual TokenType = "GreaterEqual"
	Less         TokenType = "Less"
	LessEqual    TokenType = "LessEqual"

	// Literals
	Identifier TokenType = "Identifier"
	Number     TokenType = "Number"
	String     TokenType = "String"

	// Keywords
	And     TokenType = "And"
	Assert  TokenType = "Assert"
	Break   TokenType = "Break"
	Class   TokenType = "Class"
	Echo    TokenType = "Echo"
	Else    TokenType = "Else"
	False   TokenType = "False"
	Fn      TokenType = "Fn"
	For     TokenType = "For"
	Foreign TokenType = "Foreign"
	Import  TokenType = "Import"
	Nil     TokenType = "Nil"
	Or      TokenType = "Or"
	Return  TokenType = "Return"
	True    TokenType = "True"
	Var     TokenType = "Var"
	While   TokenType = "While"

	Eof     TokenType = "Eof"
	Newline TokenType = "Newline"

	Error TokenType = "Error"
)
