package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kounoike/dtv-discord-go/db"
)

type ProgramMessageTemplateArgs struct {
	Program db.Program
	Service db.Service
}

var programMessageTemplateString = `{{ .Program.Name }}
{{ .Program.Description }}
{{ .Service.Name }}
{{ .Program.Json | toExtendStr }}
{{ .Program.StartAt |toTimeStr }}～{{ .Program.Duration | toDurationStr }}
==============================================================================================`

func toTimeStr(t int64) string {
	return time.Unix(t/1000, (t%1000)*1000).Format("2006/01/02(Mon) 03:04")
}

func toDurationStr(d int32) string {
	hour := d / (60 * 60 * 1000)
	hourStr := ""
	if hour > 0 {
		hourStr = fmt.Sprintf("%d時間", hour)
	}
	min := d % (60 * 60 * 1000) / (60 * 1000)
	minStr := fmt.Sprintf("%d分", min)
	return hourStr + minStr
}

func toExtendStr(j json.RawMessage) string {
	b, _ := j.MarshalJSON()
	any := jsoniter.Get(b, "extended")
	str := ""
	for _, key := range any.Keys() {
		str += fmt.Sprintf("%s:%s", key, any.Get(key).ToString())
	}
	return str
}

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
	args := ProgramMessageTemplateArgs{
		Program: program,
		Service: service,
	}
	err = tmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	weekdayja := strings.NewReplacer(
		"Sun", "日",
		"Mon", "月",
		"Tue", "火",
		"Wed", "水",
		"Thu", "木",
		"Fri", "金",
		"Sat", "土",
	)
	return weekdayja.Replace(b.String()), nil
}
