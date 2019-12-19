package parser

import (
	"fmt"

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
