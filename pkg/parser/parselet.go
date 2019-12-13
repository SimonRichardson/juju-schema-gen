package parser

import (
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Reader interface {
	Peek(int) (lexer.Token, error)
	Advance(int)
}

type Parselet interface {
	Parse(Reader, lexer.Token) (Expression, error)
}

type Expression interface {
	Type() ExpressionType
	Tokens() []lexer.Token
}

type ExpressionType int

const (
	EKeyword ExpressionType = iota
	EVersion ExpressionType = iota
)
