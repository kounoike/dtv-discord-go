package template

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	weekdayja = strings.NewReplacer(
		"Sun", "日",
		"Mon", "月",
		"Tue", "火",
		"Wed", "水",
		"Thu", "木",
		"Fri", "金",
		"Sat", "土",
	)
)

func toTimeStr(t int64) string {
	return time.Unix(t/1000, (t%1000)*1000).Format("2006/01/02(Mon) 15:04")
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
		str += fmt.Sprintf("%s:%s\n", key, any.Get(key).ToString())
	}
	return str
}
