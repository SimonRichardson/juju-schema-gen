package parser

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Version struct {
	Left     lexer.TokenType
	Right    lexer.TokenType
	position cursor.Position
	tokens   []lexer.Token
}

func (k *Version) Type() ExpressionType {
	return EVersion
}

func (k *Version) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Version) Position() cursor.Position {
	return k.position
}

func (k *Version) Parse(reader Reader, token lexer.Token) (Expression, error) {
	if token.Type != k.Left {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", token.Type.String(), k.Left.String())
	}
	version, err := reader.Peek(0)
	if err != nil {
		return nil, err
	}
	if version.Type != lexer.TNumber {
		return nil, fmt.Errorf("unexpected type %q, wanted \"0-9\"", version.Type.String())
	}
	right, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	if right.Type != k.Right {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", right.Type.String(), k.Right.String())
	}
	reader.AdvanceTo(2)
	return &Version{
		position: token.Position,
		tokens:   []lexer.Token{token, version, right},
	}, nil
}
