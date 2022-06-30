package logger

import (
	"fmt"
	"time"
)

//Level 日志级别
type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
)

var (
	levelTable = map[Level]string{
		Debug: "Debug",
		Info:  "Info",
		Warn:  "Warn",
		Error: "Error",
	}
)

//Logger 处理日志的logger
type Logger struct {
	name  string   //logger名称
	level Level    //日志级别
	dst   Appender //写入日志的目的地
}

func New(name string, level Level, dst Appender) *Logger {
	return &Logger{
		name:  name,
		level: level,
		dst:   dst,
	}
}

//格式化日志信息并写入
func (l *Logger) log(level Level, msg string, params ...any) {
	logMsg := fmt.Sprintf(msg, params...)
	//当前时间
	nowTime := time.Now()
	timeStr := nowTime.Format("01-02 15:04:05")
	//日志级别
	levelStr := levelTable[level]
	log := fmt.Sprintf("[%s]%s %s: %s\n", levelStr, timeStr, l.name, logMsg)
	l.dst.WriteMsg(log)
}

//Debug debug级别日志
func (l *Logger) Debug(msg string, params ...any) {
	if l.level > Debug {
		return
	}
	l.log(Debug, msg, params...)
}

//Info info级别日志
func (l *Logger) Info(msg string, params ...any) {
	if l.level > Info {
		return
	}
	l.log(Info, msg, params...)
}

//Warn warn级别日志
func (l *Logger) Warn(msg string, params ...any) {
	if l.level > Warn {
		return
	}
	l.log(Warn, msg, params...)
}

//Error error级别日志
func (l *Logger) Error(msg string, params ...any) {
	l.log(Error, msg, params...)
}
