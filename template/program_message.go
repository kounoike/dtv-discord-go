package template

import (
	"bytes"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"
)

type ProgramMessageTemplateArgs struct {
	Program db.Program
	Service mirakc_model.Service
}

var programMessageTemplateString = `==============================================================================================
{{ .Program.Name }}
{{ .Program.Description }}
{{ .Service.Name }}
`

func GetProgramMessage(program db.Program, service mirakc_model.Service) (string, error) {
	tmpl, err := template.New("program-message").Parse(programMessageTemplateString)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	args := ProgramMessageTemplateArgs{
		Program: program,
		Service: service,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
