package errors

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
)

type OverflowError struct {
	Context      string
	Char         string
	Position     cursor.Position
	Alternatives []string
}

func (e OverflowError) Error() string {
	t := template.Must(template.New("error").Funcs(funcs).Parse(overflowTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, e); err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Unexpected overflow found %q", string(e.Char))
	}
	return strings.TrimSpace(buf.String())
}

const overflowTemplate = `Unexpected overflow found "{{.Char}}"
{{$position := .Position}}
{{$position.Line}}| {{.Context}}
{{offset (len (print $position.Line)) 2}}{{underline $position.Start $position.End}}

{{if .Alternatives}}
Maybe you want one of the following?

{{range $alt := .Alternatives}}
  - {{$alt}}
{{end}}
{{end}}
`
