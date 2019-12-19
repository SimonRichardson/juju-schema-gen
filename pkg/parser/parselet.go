package parser

import (
	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Reader interface {
	Peek(int) (lexer.Token, error)
	AdvanceTo(int)
	Parse() ([]Expression, int, error)
	Len() int
}

type Parselet interface {
	Parse(Reader, lexer.Token) (Expression, error)
}

type Expression interface {
	Type() ExpressionType
	Tokens() []lexer.Token
	Position() cursor.Position
}

type RecursiveExpression interface {
	Expressions() []Expression
}

type ExpressionType int

func (t ExpressionType) String() string {
	switch t {
	case EKeyword:
		return "keyword"
	case EVersion:
		return "version"
	case EBody:
		return "body"
	case EType:
		return "type"
	case EList:
		return "list"
	default:
		return "UNKNOWN"
	}
}

const (
	EKeyword ExpressionType = iota
	EVersion
	EBody
	EType
	EList
)
