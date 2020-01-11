package compiler

import (
	"fmt"
	"github.com/adamjedlicka/go-blu/src/parser"
	"github.com/adamjedlicka/go-blu/src/value"
	"os"
	"strconv"
)

const MaxLocals = 65000

type Compiler struct {
	p     *parser.Parser
	chunk *Chunk

	locals   []Local
	upvalues []Upvalue

	scopeDepth int8

	hadError  bool
	panicMode bool
}

func NewCompiler(name string, source string) Compiler {
	return Compiler{
		p:     parser.NewParser(source),
		chunk: NewChunk(name),

		locals:   make([]Local, 0),
		upvalues: make([]Upvalue, 0),

		scopeDepth: 0,

		hadError:  false,
		panicMode: false,
	}
}

func (c *Compiler) Compile() *Chunk {
	for true {
		c.advance()

		if !c.check(parser.Newline) {
			break
		}
	}

	for !c.match(parser.Eof) {
		c.declaration()
	}

	// Patch last Pop for REPL
	if len(c.chunk.code) > 0 && c.chunk.code[len(c.chunk.code)-1] == uint8(Pop) {
		c.chunk.code[len(c.chunk.code)-1] = uint8(Return)
	} else {
		c.emitReturn()
	}

	if c.hadError {
		return nil
	}

	return c.chunk
}

func (c *Compiler) declaration() {
	if c.match(parser.Var) {
		c.varDeclaration()
	} else {
		c.statement()
	}

	if c.panicMode {
		c.synchronize()
	}
}

func (c *Compiler) varDeclaration() {
	name := c.parseVariable("Expect variable name.")

	if c.match(parser.Equal) {
		c.expression()
	} else {
		c.emitOpCode(Nil)
	}

	c.defineVariable(name)

	c.expectNewlineOrSemicolon()
}

func (c *Compiler) addLocal(name parser.Token) {
	if len(c.locals) == MaxLocals {
		c.error("Too many local variables in function.")
		return
	}

	local := Local{
		name:      name,
		depth:     -1,
		isUpvalue: false,
	}

	c.locals = append(c.locals, local)
}

func (c *Compiler) declareVariable() {
	// Global variables are implicitly declared.
	if c.scopeDepth == 0 {
		return
	}

	name := c.p.Previous()

	for i := len(c.locals) - 1; i >= 0; i-- {
		local := c.locals[i]

		if local.depth != -1 && local.depth < c.scopeDepth {
			break
		}

		if name.Lexeme() == local.name.Lexeme() {
			c.error("Variable with this name already declared in this scope.")
		}
	}

	c.addLocal(name)
}

func (c *Compiler) markInitialized() {
	if c.scopeDepth == 0 {
		return
	}

	c.locals[len(c.locals)-1].depth = c.scopeDepth
}

func (c *Compiler) defineVariable(index uint16) {
	if c.scopeDepth > 0 {
		c.markInitialized()
		return
	}

	c.emitOpCode(DefineGlobal)
	c.emitShort(index)
}

func (c *Compiler) identifierConstant(name parser.Token) uint16 {
	return c.makeConstant(value.String(name.Lexeme()))
}

func (c *Compiler) parseVariable(message string) uint16 {
	c.consume(parser.Identifier, message)

	c.declareVariable()
	if c.scopeDepth > 0 {
		return 0
	}

	return c.identifierConstant(c.p.Previous())
}

func (c *Compiler) statement() {
	if c.match(parser.LeftBrace) {
		c.beginScope()
		c.block()
		c.endScope()
	} else if c.match(parser.If) {
		c.ifStatement()
	} else if c.match(parser.Return) {
		c.returnStatement()
	} else {
		c.expressionStatement()
	}
}

func (c *Compiler) block() {
	for !c.check(parser.RightBrace) && !c.check(parser.Eof) {
		c.declaration()
	}

	c.consume(parser.RightBrace, "Expect '}' after block.")
}

func (c *Compiler) beginScope() {
	c.scopeDepth++
}

func (c *Compiler) endScope() {
	c.scopeDepth--

	for len(c.locals) > 0 && c.locals[len(c.locals)-1].depth > c.scopeDepth {
		if c.locals[len(c.locals)-1].isUpvalue {
			// TODO : Implement closing of upvalues
			panic("unimplemented")
		} else {
			c.emitOpCode(Pop)
		}

		c.locals = c.locals[:len(c.locals)-1]
	}
}

func (c *Compiler) ifStatement() {
	c.expression()
	ifJump := c.emitJump(JumpIfFalsy)
	c.emitOpCode(Pop) // Condition

	// One-line notation
	if c.match(parser.Colon) {
		c.statement()

		elseJump := c.emitJump(Jump)
		c.patchJump(ifJump)
		c.emitOpCode(Pop) // Condition
		c.patchJump(elseJump)

		return
	}

	c.consume(parser.LeftBrace, "Expect '{' after if condition.")
	c.beginScope()
	c.block()
	c.endScope()

	c.patchJump(ifJump)
	c.emitOpCode(Pop) // Condition

	if c.match(parser.Else) {
		if c.match(parser.LeftBrace) {
			c.beginScope()
			c.block()
			c.endScope()
		} else if c.match(parser.If) {
			c.ifStatement()
		} else {
			c.error("Expect 'if' or '{' after 'else'.")
		}
	}
}

func (c *Compiler) returnStatement() {
	if c.match(parser.Newline) {
		c.emitReturn()
	} else {
		needsNewline := !c.check(parser.Fn)

		c.expression()
		c.emitOpCode(Return)

		if needsNewline {
			c.expectNewlineOrSemicolon()
		}
	}
}

func (c *Compiler) expressionStatement() {
	c.expression()

	c.emitOpCode(Pop)

	c.expectNewlineOrSemicolon()
}

func (c *Compiler) expression() {
	c.parsePrecedence(PrecedenceAssignment)
}

func (c *Compiler) parsePrecedence(precedence Precedence) {
	c.advance()

	prefixRule := parseRules[c.p.Previous().Type()].prefix
	if prefixRule == nil {
		c.error("Expect expression.")

		return
	}

	canAssign := precedence <= PrecedenceAssignment
	prefixRule(c, canAssign)

	for precedence <= parseRules[c.p.Current().Type()].precedence {
		c.advance()

		infixRule := parseRules[c.p.Previous().Type()].infix
		infixRule(c, canAssign)
	}

	if canAssign && c.match(parser.Equal) {
		c.error("Invalid assignment target.")

		// Parse the expression so compiler prints proper error messages.
		c.expression()
	}
}

func (c *Compiler) unary(canAssign bool) {
	operatorType := c.p.Previous().Type()

	c.parsePrecedence(PrecedenceUnary)

	switch operatorType {
	case parser.Bang:
		c.emitOpCode(Not)
	case parser.Minus:
		c.emitOpCode(Negate)
	default:
		panic("unreachable")
	}
}

func (c *Compiler) binary(canAssign bool) {
	operatorType := c.p.Previous().Type()

	rule := parseRules[operatorType]

	c.parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case parser.EqualEqual:
		c.emitOpCode(Equal)
	case parser.Greater:
		c.emitOpCode(Greater)
	case parser.GreaterEqual:
		c.emitOpCode(GreaterEqual)
	case parser.Less:
		c.emitOpCode(Less)
	case parser.LessEqual:
		c.emitOpCode(LessEqual)
	case parser.BangEqual:
		c.emitOpCode(NotEqual)

	case parser.Plus:
		c.emitOpCode(Add)
	case parser.Minus:
		c.emitOpCode(Subtract)
	case parser.Slash:
		c.emitOpCode(Divide)
	case parser.Star:
		c.emitOpCode(Multiply)
	case parser.Caret:
		c.emitOpCode(Exponentiate)
	case parser.Percent:
		c.emitOpCode(Reminder)
	}
}

func (c *Compiler) number(canAssign bool) {
	lexeme := c.p.Previous().Lexeme()
	number, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		panic(err)
	}

	c.emitConstant(value.Number(number))
}

func (c *Compiler) string(canAssign bool) {
	lexeme := c.p.Previous().Lexeme()
	string := lexeme[1 : len(lexeme)-1]

	c.emitConstant(value.String(string))
}

func (c *Compiler) literal(canAssign bool) {
	switch c.p.Previous().Type() {
	case parser.False:
		c.emitOpCode(False)
	case parser.True:
		c.emitOpCode(True)
	case parser.Nil:
		c.emitOpCode(Nil)
	default:
		panic("unreachable")
	}
}

func (c *Compiler) grouping(canAssign bool) {
	c.expression()
	c.consume(parser.RightParen, "Expect ')' after expression")
}

func (c *Compiler) resolveLocal(name parser.Token) (uint16, bool) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		local := c.locals[i]
		if local.name.Lexeme() == name.Lexeme() {
			if local.depth == -1 {
				c.error("Cannot read local variable in its own initializer.")
			}

			return uint16(i), true
		}
	}

	return 0, false
}

func (c *Compiler) namedVariable(name parser.Token, canAssign bool) {
	var getOp OpCode
	var setOp OpCode

	arg, ok := c.resolveLocal(name)

	if ok {
		getOp = GetLocal
		setOp = SetLocal
	} else {
		arg = c.identifierConstant(name)
		getOp = GetGlobal
		setOp = SetGlobal
	}

	if canAssign && c.match(parser.Equal) {
		c.expression()
		c.emitOpCode(setOp)
		c.emitShort(arg)
	} else {
		c.emitOpCode(getOp)
		c.emitShort(arg)
	}
}

func (c *Compiler) variable(canAssign bool) {
	c.namedVariable(c.p.Previous(), canAssign)
}

func (c *Compiler) emitByte(byte uint8) {
	c.chunk.pushCode(byte, c.p.Current().Line())
}

func (c *Compiler) emitShort(short uint16) {
	c.chunk.pushCode(uint8((short>>8)&0xff), c.p.Current().Line())
	c.chunk.pushCode(uint8(short&0xff), c.p.Current().Line())
}

func (c *Compiler) emitOpCode(opCode OpCode) {
	c.chunk.pushCode(uint8(opCode), c.p.Current().Line())
}

func (c *Compiler) emitJump(code OpCode) int {
	c.emitOpCode(code)
	c.emitShort(0)

	return len(c.chunk.code) - 2
}

func (c *Compiler) patchJump(jump int) {
	length := len(c.chunk.code) - 2 - jump

	c.chunk.code[jump] = uint8((length >> 8) & 0xff)
	c.chunk.code[jump+1] = uint8(length & 0xff)
}

func (c *Compiler) makeConstant(value value.Value) uint16 {
	constant := c.chunk.pushConstant(value)
	if constant > MaxConstants {
		c.error("Too many constants in one chunk.")
		return 0
	}

	return constant
}

func (c *Compiler) emitConstant(value value.Value) {
	constant := c.makeConstant(value)

	c.emitOpCode(Constant)
	c.emitShort(constant)
}

func (c *Compiler) emitReturn() {
	c.emitOpCode(Nil)
	c.emitOpCode(Return)
}

func (c *Compiler) consumeNewlines() {
	for c.p.Current().Type() == parser.Newline {
		c.p.SetCurrent(c.p.NextToken())
	}
}

func (c *Compiler) skipNewlines() {
	switch c.p.Previous().Type() {
	case parser.Newline, parser.LeftBrace, parser.RightBrace, parser.Semicolon, parser.Dot:
		c.consumeNewlines()
	}
}

func (c *Compiler) expectNewlineOrSemicolon() {
	// TODO : Improve consuming of newlines
	// If previous token is RightBrace then all newlines were already consumed
	if c.p.Previous().Type() == parser.RightBrace {
		return
	}

	// If current token is RightBrace then we don't need newline nor semicolon
	if c.p.Current().Type() == parser.RightBrace {
		return
	}

	// If we are at the end of the file then we don't need newline nor semicolon
	if c.p.Previous().Type() == parser.Eof || c.p.Current().Type() == parser.Eof {
		return
	}

	if !c.match(parser.Semicolon) {
		c.consume(parser.Newline, "Expect newline or ';'.")
	}
}

func (c *Compiler) advance() {
	c.p.SetPrevious(c.p.Current())

	for true {
		c.p.SetCurrent(c.p.NextToken())
		if c.p.Current().Type() != parser.Error {
			break
		}

		c.errorAtCurrent(c.p.Current().Lexeme())
	}

	c.skipNewlines()
}

// Checks whether next token is of the given type.
// Returns true if so, otherwise returns false.
func (c *Compiler) check(tokenType parser.TokenType) bool {
	if tokenType == parser.Newline && c.p.Previous().Type() == parser.RightBrace {
		return true
	}

	return c.p.Current().Type() == tokenType
}

func (c *Compiler) consume(tokenType parser.TokenType, message string) {
	if c.check(tokenType) {
		c.advance()

		return
	}

	c.errorAtCurrent(message)
}

// Checks whether next token is of the given type.
// If yes, consumes it and returns true, otherwise it does not consume any tokens and return false.
func (c *Compiler) match(tokenType parser.TokenType) bool {
	if !c.check(tokenType) {
		return false
	}

	c.advance()

	return true
}

func (c *Compiler) synchronize() {
	c.panicMode = false

	for c.p.Current().Type() != parser.Eof {
		switch c.p.Current().Type() {
		case parser.Class, parser.Fn, parser.Var, parser.For, parser.If, parser.While:
			return
		default:
			c.advance()
		}
	}
}

func (c *Compiler) error(message string) {
	c.errorAt(c.p.Previous(), message)
}

func (c *Compiler) errorAtCurrent(message string) {
	c.errorAt(c.p.Current(), message)
}

func (c *Compiler) errorAt(token parser.Token, message string) {
	if c.panicMode {
		return
	}

	c.panicMode = true

	_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line())

	switch token.Type() {
	case parser.Eof:
		_, _ = fmt.Fprintf(os.Stderr, " at end")
	case parser.Newline:
		_, _ = fmt.Fprintf(os.Stderr, " at newline")
	default:
		_, _ = fmt.Fprintf(os.Stderr, " at '%s'", token.Lexeme())
	}

	_, _ = fmt.Fprintf(os.Stderr, ": %s\n", message)

	c.hadError = true
}
