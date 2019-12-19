package ast

import (
	"fmt"
	"io"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type Data struct {
	cursor.Position
	Name   parser.Expression
	Body   parser.Expression
	Fields []Field
}

func (d *Data) Read(expressions []parser.Expression) (int, error) {
	exp, err := matchString(expressions, 0, "data")
	if err != nil {
		return -1, err
	}
	d.Position = exp.Position()

	if d.Name, _, err = pluckString(expressions, 1); err != nil {
		return -1, err
	}
	if d.Body, err = pluckBody(expressions, 2); err != nil {
		return -1, err
	}

	return 2, nil
}

func (d *Data) Format(out io.Writer) error {
	fmt.Fprintf(out, "data %s %s\n", tokenStrings(d.Name.Tokens()), firstTokenString(d.Body.Tokens()))
	fmt.Fprintf(out, "%s\n", lastTokenString(d.Body.Tokens()))
	return nil
}
