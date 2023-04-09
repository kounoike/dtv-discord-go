package template

import (
	"bytes"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
)

type recordingStartedMessageTemplateArgs struct {
	Program     db.Program
	Service     db.Service
	ContentPath string
}

const (
	recordingStartedMessageTemplateString = `**録画開始**：{{ .Program.Name }}
　{{ .Service.Name }} {{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
　保存先：` + "`" + `{{ .ContentPath }}` + "`"
)

func GetRecordingStartedMessage(program db.Program, service db.Service, contentPath string) (string, error) {
	funcMap := map[string]interface{}{
		"toTimeStr":     toTimeStr,
		"toDurationStr": toDurationStr,
		"toExtendStr":   toExtendStr,
	}
	tmpl, err := template.New("recording-started-message").Funcs(funcMap).Parse(recordingStartedMessageTemplateString)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	args := recordingStartedMessageTemplateArgs{
		Program:     program,
		Service:     service,
		ContentPath: contentPath,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return weekdayja.Replace(b.String()), nil
}
