package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"github.com/jellycheng/gosupport"
)

var globalCfg *gosupport.DataManage
func init()  {
	globalCfg = gosupport.NewGlobalCfgSingleton()
}

type LogFormatter struct {
	log.Formatter
}

/**
【格式】日期+时间+毫秒 | 服务名  | 日志级别 |traceid | 日志内容 | 日志的详细参数json格式
【示例】 2019-02-28 11:23:33 998 user-service 10000 DEBUG tr15513242132905 "调试信息"  {"1":"男","2":”女","3":"保密"}
*/
func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var buffer = &bytes.Buffer{}
	buffer.WriteString(entry.Time.Format("2006-01-02 15:04:05.000"))
	buffer.WriteString(" ")
	buffer.WriteString(globalCfg.GetString("app_name"))
	buffer.WriteString(" ")

	buffer.WriteString(strings.ToUpper(entry.Level.String()))
	buffer.WriteString(" ")

	traceid:=entry.Data["traceid"]
	if v,ok:=traceid.(string);ok {
		buffer.WriteString(v)
		buffer.WriteString(" ")
	}

	buffer.WriteString(entry.Message)
	buffer.WriteString(" ")

	if len(entry.Data) > 0 {
		jsonData := ""
		if d, err:= json.Marshal(entry.Data); err==nil {
			jsonData = string(d)
		}
		buffer.WriteString(fmt.Sprintf(" %v", jsonData))
	}
	buffer.WriteString("\n")

	return buffer.Bytes(), nil
}


