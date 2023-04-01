package template

import (
	"bytes"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
)

type recordingFailedMessageTemplateArgs struct {
	Program     db.Program
	Service     db.Service
	ContentPath string
	Reason      string
}

const (
	recordingFailedMessageTemplateString = `**録画失敗**：{{ .Program.Name }}
{{ .Service.Name }} {{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
保存先：` + "`" + `{{ .ContentPath }}` + "`" + `
エラーメッセージ：{{ .Reason }}`
)

func GetRecordingFailedMessage(program db.Program, service db.Service, contentPath string, reason string) (string, error) {
	funcMap := map[string]interface{}{
		"toTimeStr":     toTimeStr,
		"toDurationStr": toDurationStr,
		"toExtendStr":   toExtendStr,
	}
	tmpl, err := template.New("recording-failed-message").Funcs(funcMap).Parse(recordingFailedMessageTemplateString)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	args := recordingFailedMessageTemplateArgs{
		Program:     program,
		Service:     service,
		ContentPath: contentPath,
		Reason:      reason,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return weekdayja.Replace(b.String()), nil
}
