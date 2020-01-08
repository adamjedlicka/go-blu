package parser

var keywords = map[string]TokenType{
	"and":     And,
	"assert":  Assert,
	"break":   Break,
	"class":   Class,
	"echo":    Echo,
	"else":    Else,
	"false":   False,
	"fn":      Fn,
	"for":     For,
	"Foreign": Foreign,
	"import":  Import,
	"nil":     Nil,
	"or":      Or,
	"return":  Return,
	"true":    True,
	"var":     Var,
	"while":   While,
}
