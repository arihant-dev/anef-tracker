package log

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Level string
type Format string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"

	FormatText Format = "text"
	FormatJSON Format = "json"
)

var currentFormat Format = FormatText

func SetFormat(fmtStr string) {
	if strings.ToLower(fmtStr) == "json" {
		currentFormat = FormatJSON
	} else {
		currentFormat = FormatText
	}
}

func Info(msg string, keyvals ...interface{}) {
	logMessage(LevelInfo, msg, keyvals...)
}

func Warn(msg string, keyvals ...interface{}) {
	logMessage(LevelWarn, msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	logMessage(LevelError, msg, keyvals...)
}

func Debug(msg string, keyvals ...interface{}) {
	logMessage(LevelDebug, msg, keyvals...)
}

func logMessage(lvl Level, msg string, keyvals ...interface{}) {
	now := time.Now()
	kvMap := make(map[string]interface{})
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			kStr := fmt.Sprintf("%v", keyvals[i])
			kvMap[kStr] = keyvals[i+1]
		}
	}

	if currentFormat == FormatJSON {
		kvMap["time"] = now.Format(time.RFC3339)
		kvMap["level"] = string(lvl)
		kvMap["msg"] = msg
		b, _ := json.Marshal(kvMap)
		log.Println(string(b))
		return
	}

	var kvs []string
	for k, v := range kvMap {
		kvs = append(kvs, fmt.Sprintf("%s=%v", k, v))
	}

	kvStr := ""
	if len(kvs) > 0 {
		kvStr = " | " + strings.Join(kvs, " ")
	}

	log.Printf("[%s] [%s] %s%s", now.Format("15:04:05"), lvl, msg, kvStr)
}
