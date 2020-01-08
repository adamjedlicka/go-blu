package parser

import (
	"testing"
)

func TestItParsesBasicTokens(t *testing.T) {
	p := NewParser("+()>>=<")

	if p.NextToken().tokenType != Plus {
		t.Error("Expect 'Plus'")
	}

	if p.NextToken().tokenType != LeftParen {
		t.Error("Expect 'LeftParen'")
	}

	if p.NextToken().tokenType != RightParen {
		t.Error("Expect 'RightParen'")
	}

	if p.NextToken().tokenType != Greater {
		t.Error("Expect 'Greater'")
	}

	if p.NextToken().tokenType != GreaterEqual {
		t.Error("Expect 'GreaterEqual'")
	}

	if p.NextToken().tokenType != Less {
		t.Error("Expect 'Less'")
	}

	if p.NextToken().tokenType != Eof {
		t.Error("Expect 'Eof'")
	}
}

func TestItParsesNumbers(t *testing.T) {
	p := NewParser("1 10 1.1 999.83")

	if p.NextToken().lexeme != "1" {
		t.Error("Expect '1'")
	}

	if p.NextToken().lexeme != "10" {
		t.Error("Expect '10'")
	}

	if token := p.NextToken(); token.lexeme != "1.1" {
		t.Error(token)
	}

	if token := p.NextToken(); token.lexeme != "999.83" {
		t.Error(token)
	}
}

func TestItParsesStrings(t *testing.T) {
	p := NewParser("\"Hello, World!\" \"WHAT\" \"\"")

	if p.NextToken().lexeme != "\"Hello, World!\"" {
		t.Error("Expect '\"Hello, World!\"'")
	}

	if p.NextToken().lexeme != "\"WHAT\"" {
		t.Error("Expect '\"WHAT\"'")
	}

	if p.NextToken().tokenType != String {
		t.Error("Expect 'String'")
	}
}

func TestItParsesKeywords(t *testing.T) {
	p := NewParser("and or class fn for while")

	if p.NextToken().tokenType != And {
		t.Error("Expect 'And'")
	}

	if p.NextToken().tokenType != Or {
		t.Error("Expect 'Or'")
	}

	if p.NextToken().tokenType != Class {
		t.Error("Expect 'Class'")
	}

	if p.NextToken().tokenType != Fn {
		t.Error("Expect 'Fn'")
	}

	if p.NextToken().tokenType != For {
		t.Error("Expect 'For'")
	}

	if p.NextToken().tokenType != While {
		t.Error("Expect 'While'")
	}
}

func TestItParsesIdentifiers(t *testing.T) {
	p := NewParser("Class myClass WHAT _leading _1")

	if p.NextToken().tokenType != Identifier {
		t.Error("Expect 'Identifier'")
	}

	if p.NextToken().tokenType != Identifier {
		t.Error("Expect 'Identifier'")
	}

	if p.NextToken().tokenType != Identifier {
		t.Error("Expect 'Identifier'")
	}

	if p.NextToken().tokenType != Identifier {
		t.Error("Expect 'Identifier'")
	}

	if p.NextToken().lexeme != "_1" {
		t.Error("Expect '_1'")
	}
}

func TestItSkipsWhitespace(t *testing.T) {
	p := NewParser("   true   false   ")

	if token := p.NextToken(); token.tokenType != True {
		t.Error(token)
	}

	if token := p.NextToken(); token.tokenType != False {
		t.Error(token)
	}

	if token := p.NextToken(); token.tokenType != Eof {
		t.Error(token)
	}
}
