package ast

import (
	"bytes"
	"io"

	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type AST struct {
	Package *Package
}

func (a AST) Format(out io.Writer) error {
	if err := a.Package.Format(out); err != nil {
		return err
	}
	return nil
}

func (a AST) String() string {
	var b bytes.Buffer
	_ = a.Format(&b)
	return b.String()
}

func Generate(expressions []parser.Expression) (AST, error) {
	pack := &Package{}
	_, err := pack.Read(expressions[0:])
	if err != nil {
		return AST{}, err
	}
	return AST{
		Package: pack,
	}, nil
}
