package parser

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Type struct {
	Left     lexer.TokenType
	Right    lexer.TokenType
	Keyword  lexer.TokenType
	position cursor.Position
	tokens   []lexer.Token
}

func (k *Type) Type() ExpressionType {
	return EType
}

func (k *Type) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Type) Position() cursor.Position {
	return k.position
}

func (k *Type) Parse(reader Reader, token lexer.Token) (Expression, error) {
	if token.Type != k.Left {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", token.Type.String(), k.Left.String())
	}
	right, err := reader.Peek(0)
	if err != nil {
		return nil, err
	}
	if right.Type != k.Right {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", right.Type.String(), k.Right.String())
	}
	keyword, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	if keyword.Type != k.Keyword {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", keyword.Type.String(), k.Keyword.String())
	}
	reader.AdvanceTo(2)
	return &Type{
		position: token.Position,
		tokens:   []lexer.Token{token, right, keyword},
	}, nil
}
