package parser

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Body struct {
	Left        lexer.TokenType
	Right       lexer.TokenType
	position    cursor.Position
	tokens      []lexer.Token
	expressions []Expression
}

func (k *Body) Type() ExpressionType {
	return EBody
}

func (k *Body) Tokens() []lexer.Token {
	return k.tokens
}

func (k *Body) Position() cursor.Position {
	return k.position
}

func (k *Body) Expressions() []Expression {
	return k.expressions
}

func (k *Body) Parse(reader Reader, token lexer.Token) (Expression, error) {
	if token.Type != k.Left {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", token.Type.String(), k.Left.String())
	}

	tokens := []lexer.Token{
		token,
	}
	expressions := make([]Expression, 0)
	for i := 0; i < reader.Len(); i++ {
		// Attempt to read in the close, if it matches the break.
		right, err := reader.Peek(i)
		if err != nil {
			return nil, err
		}
		if right.Type == k.Right {
			reader.AdvanceTo(i + 1)
			token, err := k.Right.Token()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			break
		}

		// Parse recursively
		parsed, consumed, err := reader.Parse()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, parsed...)
		i += consumed
	}

	return &Body{
		position:    token.Position,
		tokens:      tokens,
		expressions: expressions,
	}, nil
}
