package parser

import (
	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Keyword struct {
	position cursor.Position
	tokens   []lexer.Token
}

func (k *Keyword) Type() ExpressionType {
	return EKeyword
}

func (k *Keyword) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Keyword) Position() cursor.Position {
	return k.position
}

func (k *Keyword) Parse(reader Reader, token lexer.Token) (Expression, error) {
	return &Keyword{
		position: token.Position,
		tokens:   []lexer.Token{token},
	}, nil
}
