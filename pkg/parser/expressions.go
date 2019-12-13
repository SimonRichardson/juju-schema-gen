package parser

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Keyword struct {
	tokens []lexer.Token
}

func (k *Keyword) Type() ExpressionType {
	return EKeyword
}

func (k *Keyword) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Keyword) Parse(reader Reader, token lexer.Token) (Expression, error) {
	return &Keyword{
		tokens: []lexer.Token{token},
	}, nil
}

type Version struct {
	tokens []lexer.Token
}

func (k *Version) Type() ExpressionType {
	return EVersion
}

func (k *Version) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Version) Parse(reader Reader, token lexer.Token) (Expression, error) {
	if token.Type != lexer.TLeftSquareBracket {
		return nil, fmt.Errorf("unexpected type %v, wanted '['", token.Type)
	}
	version, err := reader.Peek(0)
	if err != nil {
		return nil, err
	}
	if version.Type != lexer.TNumber {
		return nil, fmt.Errorf("unexpected type %v, wanted '0-9'", version.Type)
	}
	right, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	if right.Type != lexer.TRightSquareBracket {
		return nil, fmt.Errorf("unexpected type %v, wanted ']'", right.Type)
	}
	reader.Advance(2)
	return &Version{
		tokens: []lexer.Token{token, version, right},
	}, nil
}
