package template

import (
	"bytes"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
)

type programMessageTemplateArgs struct {
	Program db.Program
	Service db.Service
}

const (
	programMessageTemplateString = `{{ .Program.Name }}
{{ .Program.Genre }}
{{ .Program.Description }}
{{ .Service.Name }}
{{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
{{ .Program.Json | toExtendStr }}`
)

func GetProgramMessage(program db.Program, service db.Service) (string, error) {
	funcMap := map[string]interface{}{
		"toTimeStr":     toTimeStr,
		"toDurationStr": toDurationStr,
		"toExtendStr":   toExtendStr,
	}
	tmpl, err := template.New("program-message").Funcs(funcMap).Parse(programMessageTemplateString)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	args := programMessageTemplateArgs{
		Program: program,
		Service: service,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return weekdayja.Replace(b.String()), nil
}
