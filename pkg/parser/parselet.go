package parser

import (
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Reader interface {
	Peek(int) (lexer.Token, error)
	AdvanceTo(int)
	Parse() ([]Expression, error)
}

type Parselet interface {
	Parse(Reader, lexer.Token) (Expression, error)
}

type Expression interface {
	Type() ExpressionType
	Tokens() []lexer.Token
}

type RecursiveExpression interface {
	Expressions() []Expression
}

type ExpressionType int

const (
	EKeyword ExpressionType = iota
	EVersion
	EBody
	EType
)
