package errors

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type CharPositionError struct {
	Context       string
	Char          string
	Line          int
	PositionStart int
	PositionEnd   int
	Alternatives  []string
}

func (e CharPositionError) Error() string {
	t := template.Must(template.New("error").Funcs(funcs).Parse(charPositionTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, e); err != nil {
		return fmt.Sprintf("Unexpected character found %q", string(e.Char))
	}
	return strings.TrimSpace(buf.String())
}

var funcs = template.FuncMap{
	"underline": func(a, b int) string {
		return fmt.Sprintf("%s%s", strings.Repeat(" ", a), strings.Repeat("^", b-a))
	},
	"offset": func(a, b int) string {
		return strings.Repeat(" ", a+b)
	},
}

const charPositionTemplate = `Unexpected character found "{{.Char | html}}"

{{.Line}}| {{.Context | html}}
{{offset (len (print .Line)) 2}}{{underline .PositionStart .PositionEnd}}

{{if .Alternatives}}
Maybe you want one of the following?

{{range $alt := .Alternatives}}
  - {{$alt}}
{{end}}
{{end}}
`
