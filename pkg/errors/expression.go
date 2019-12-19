package errors

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
)

type ExpressionError struct {
	Context      string
	Token        string
	Position     cursor.Position
	Alternatives []string
}

func (e ExpressionError) Error() string {
	t := template.Must(template.New("error").Funcs(funcs).Parse(expressionTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, e); err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Unexpected expression found %q", string(e.Token))
	}
	return strings.TrimSpace(buf.String())
}

const expressionTemplate = `Unexpected expression found "{{.Token | html}}"
{{$position := .Position}}
{{$position.Line}}| {{.Context | html}}
{{offset (len (print $position.Line)) 2}}{{underline $position.Start $position.End}}
{{if .Alternatives -}}
Maybe you want one of the following?
{{- range $alt := .Alternatives}}
  - {{$alt}}
{{- end}}
{{- end}}
`
