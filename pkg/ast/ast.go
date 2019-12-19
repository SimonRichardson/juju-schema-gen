package ast

import (
	"bytes"
	"fmt"
	"io"

	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type AST struct {
	Package *Package
	Facades []*Facade
}

func (a AST) Format(out io.Writer) error {
	if err := a.Package.Format(out); err != nil {
		return err
	}
	for _, v := range a.Facades {
		fmt.Fprintln(out, "")
		if err := v.Format(out); err != nil {
			return err
		}
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
	ptr, err := pack.Read(expressions[0:])
	if err != nil {
		return AST{}, err
	}

	var facades []*Facade
	for i := ptr; i < len(expressions); i++ {
		facade := &Facade{}
		offset, err := facade.Read(expressions[i:])
		if err != nil {
			return AST{}, err
		}
		facades = append(facades, facade)
		i += offset

	}

	return AST{
		Package: pack,
		Facades: facades,
	}, nil
}
