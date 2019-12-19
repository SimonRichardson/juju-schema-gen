package errors

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
)

type CharPositionError struct {
	Context      string
	Char         string
	Position     cursor.Position
	Alternatives []string
}

func (e CharPositionError) Error() string {
	t := template.Must(template.New("error").Funcs(funcs).Parse(charPositionTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, e); err != nil {
		fmt.Println(err)
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

const charPositionTemplate = `Unexpected character found "{{.Char}}"
{{$position := .Position}}
{{$position.Line}}| {{.Context}}
{{offset (len (print $position.Line)) 2}}{{underline $position.Start $position.End}}
{{if .Alternatives -}}
Maybe you want one of the following?
{{- range $alt := .Alternatives}}
  - {{$alt}}
{{- end}}
{{- end}}
`
