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
	programMessageTemplateString = `**{{ .Program.Name }}**
**ジャンル**:{{ .Program.Genre }}
**説明**:{{ .Program.Description }}
**放送局**:{{ .Service.Name }} **放送時間**:{{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
**番組詳細**
{{ .Program.Json | toExtendStr }}`
)

func GetProgramMessage(program db.Program, service db.Service) (string, error) {
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
