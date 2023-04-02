package template

import (
	"bytes"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
)

type autoSearchMessageTemplateArgs struct {
	Program           db.Program
	Service           db.Service
	ProgramMessageURL string
}

var autoSearchMessageTemplateString = `{{ .Program.Name }}
{{ .Service.Name }}
{{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
{{ .ProgramMessageURL }}
`

func GetAutoSearchMessage(program db.Program, service db.Service, programMessageURL string) (string, error) {
	tmpl, err := template.New("autosearch-message").Funcs(funcMap).Parse(autoSearchMessageTemplateString)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	args := autoSearchMessageTemplateArgs{
		Program:           program,
		Service:           service,
		ProgramMessageURL: programMessageURL,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return weekdayja.Replace(b.String()), nil
}
