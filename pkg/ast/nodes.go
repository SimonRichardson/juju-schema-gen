package ast

import (
	"fmt"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/errors"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type Field struct {
	cursor.Position
	Name parser.Expression
	Type parser.Expression
}

type Methods struct {
	cursor.Position
	Name   string
	Inputs []parser.Expression
	Output []parser.Expression
}

func matchString(expressions []parser.Expression, index int, value string) (parser.Expression, error) {
	expression, token, err := pluckString(expressions, index)
	if err != nil {
		return nil, err
	}
	if !token.MatchString(value) {
		position := token.Position
		return nil, errors.ExpressionError{
			Context:      contextForLine(expressions, position.Line),
			Token:        string(token.Bytes),
			Alternatives: []string{value},
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start,
				End:   position.Start + len(token.Bytes),
			},
		}
	}
	return expression, nil
}

func pluckString(expressions []parser.Expression, index int) (parser.Expression, lexer.Token, error) {
	if index >= len(expressions) {
		return nil, lexer.Token{}, errors.OverflowError{}
	}

	expression := expressions[index]
	if expression.Type() != parser.EKeyword {
		position := expression.Position()
		tokens := tokenStrings(expression.Tokens())
		return nil, lexer.Token{}, errors.ExpressionError{
			Context:      contextForLine(expressions, position.Line),
			Token:        tokens,
			Alternatives: []string{"<string> where the string matches a-zA-Z"},
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start,
				End:   position.Start + len(tokens),
			},
		}
	}

	tokens := expression.Tokens()
	if len(tokens) != 1 {
		return nil, lexer.Token{}, errors.ExpressionError{}
	}
	return expression, tokens[0], nil
}

func tokenStrings(tokens []lexer.Token) string {
	var result string
	for _, v := range tokens {
		result += string(v.Bytes)
	}
	return result
}

func firstTokenString(tokens []lexer.Token) string {
	return string(tokens[0].Bytes)
}

func lastTokenString(tokens []lexer.Token) string {
	return string(tokens[len(tokens)-1].Bytes)
}

func contextForLine(expressions []parser.Expression, line int) string {
	var lines []parser.Expression
	for _, v := range expressions {
		if v.Position().Line == line {
			lines = append(lines, v)
		}
	}
	var result string
	for _, v := range lines {
		result += fmt.Sprintf("%s ", tokenStrings(v.Tokens()))
	}
	return result
}
