package parser

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type List struct {
	Left      lexer.TokenType
	Right     lexer.TokenType
	Separator lexer.TokenType
	position  cursor.Position
	tokens    []lexer.Token
}

func (k *List) Type() ExpressionType {
	return EList
}

func (k *List) Tokens() []lexer.Token {
	return k.tokens
}

func (k *List) Position() cursor.Position {
	return k.position
}

func (k *List) Parse(reader Reader, token lexer.Token) (Expression, error) {
	if token.Type != k.Left {
		return nil, fmt.Errorf("unexpected type %q, wanted %q", token.Type.String(), k.Left.String())
	}

	tokens := []lexer.Token{
		token,
	}
	for i := 0; ; i++ {
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
		} else if right.Type == k.Separator {
			reader.AdvanceTo(i + 1)
			continue
		}

		tokens = append(tokens, right)
	}

	return &List{
		position: token.Position,
		tokens:   tokens,
	}, nil
}
