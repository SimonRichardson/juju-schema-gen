package ast

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/errors"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type Facade struct {
	cursor.Position
	Name    parser.Expression
	Version parser.Expression
	Body    parser.Expression
	Params  []*Data
	Methods []Methods
}

func (p *Facade) Read(expressions []parser.Expression) (int, error) {
	exp, err := matchString(expressions, 0, "facade")
	if err != nil {
		return -1, err
	}
	p.Position = exp.Position()

	if p.Name, _, err = pluckString(expressions, 1); err != nil {
		return -1, err
	}
	if p.Version, _, err = pluckVersion(expressions, 2); err != nil {
		return -1, err
	}
	if p.Body, err = pluckBody(expressions, 3); err != nil {
		return -1, err
	}

	body, ok := p.Body.(parser.RecursiveExpression)
	if !ok {
		return -1, errors.ExpressionError{}
	}

	if p.Params, err = pluckParams(body.Expressions()); err != nil {
		return -1, err
	}

	return len(expressions) + 1, nil
}

func (p *Facade) Format(out io.Writer) error {
	fmt.Fprintf(out, "facade %s%s %s\n\n", tokenStrings(p.Name.Tokens()), tokenStrings(p.Version.Tokens()), firstTokenString(p.Body.Tokens()))
	for _, v := range p.Params {
		buf := new(bytes.Buffer)
		if err := v.Format(buf); err != nil {
			return err
		}
		for _, v := range strings.Split(buf.String(), "\n") {
			fmt.Fprintf(out, "    %s\n", v)
		}
	}
	fmt.Fprintf(out, "%s\n", lastTokenString(p.Body.Tokens()))
	return nil
}

func pluckVersion(expressions []parser.Expression, index int) (parser.Expression, lexer.Token, error) {
	if index >= len(expressions) {
		return nil, lexer.Token{}, errors.OverflowError{}
	}

	expression := expressions[index]
	if expression.Type() != parser.EVersion {
		position := expression.Position()
		tokens := tokenStrings(expression.Tokens())
		return nil, lexer.Token{}, errors.ExpressionError{
			Context:      contextForLine(expressions, position.Line),
			Token:        tokens,
			Alternatives: []string{"<version> where version is a number"},
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start,
				End:   position.Start + len(tokens),
			},
		}
	}

	tokens := expression.Tokens()
	if len(tokens) != 3 {
		position := expression.Position()
		tokens := tokenStrings(expression.Tokens())
		return nil, lexer.Token{}, errors.ExpressionError{
			Context:      contextForLine(expressions, position.Line),
			Token:        tokens,
			Alternatives: []string{"<version> where version is a number"},
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start + 1,
				End:   position.Start + 1 + len(tokens),
			},
		}
	}
	return expression, tokens[1], nil
}

func pluckBody(expressions []parser.Expression, index int) (parser.Expression, error) {
	if index >= len(expressions) {
		return nil, errors.OverflowError{}
	}

	expression := expressions[index]
	if expression.Type() != parser.EBody {
		position := expression.Position()
		tokens := tokenStrings(expression.Tokens())
		return nil, errors.ExpressionError{
			Context: contextForLine(expressions, position.Line),
			Token:   tokens,
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start,
				End:   position.Start + len(tokens),
			},
		}
	}

	tokens := expression.Tokens()
	if len(tokens) != 2 {
		position := expression.Position()
		tokens := tokenStrings(expression.Tokens())
		return nil, errors.ExpressionError{
			Context: contextForLine(expressions, position.Line),
			Token:   tokens,
			Position: cursor.Position{
				Line:  position.Line,
				Start: position.Start + 1,
				End:   position.Start + 1 + len(tokens),
			},
		}
	}
	return expression, nil
}

func pluckParams(expressions []parser.Expression) ([]*Data, error) {
	var params []*Data
	for i := 0; i < len(expressions); i++ {
		_, err := matchString(expressions, i, "data")
		if err != nil {
			continue
		}
		data := &Data{}
		ptr, err := data.Read(expressions[i:])
		if err != nil {
			return nil, err
		}
		params = append(params, data)
		i += ptr
	}
	return params, nil
}
