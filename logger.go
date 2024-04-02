package ezmq

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type printLogger struct {
}

func (l *printLogger) Debug(v ...interface{}) {
	nv := append([]interface{}{debugLevel, time.Now().Format(time.DateTime), fileLine()}, v...)
	fmt.Println(nv...)
}

func (l *printLogger) Info(v ...interface{}) {
	nv := append([]interface{}{infoLevel, time.Now().Format(time.DateTime), fileLine()}, v...)
	fmt.Println(nv...)
}

func (l *printLogger) Warn(v ...interface{}) {
	nv := append([]interface{}{warnLevel, time.Now().Format(time.DateTime), fileLine()}, v...)
	fmt.Println(nv...)
}

func (l *printLogger) Error(v ...interface{}) {
	nv := append([]interface{}{errorLevel, time.Now().Format(time.DateTime), fileLine()}, v...)
	fmt.Println(nv...)
}

func (l *printLogger) Debugf(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s %s %s", debugLevel, time.Now().Format(time.DateTime), fileLine(), f), args...)
}

func (l *printLogger) Infof(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s %s %s", infoLevel, time.Now().Format(time.DateTime), fileLine(), f), args...)
}

func (l *printLogger) Warnf(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s %s %s", warnLevel, time.Now().Format(time.DateTime), fileLine(), f), args...)
}

func (l *printLogger) Errorf(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s %s %s", errorLevel, time.Now().Format(time.DateTime), fileLine(), f), args...)
}

const (
	debugLevel = "DBUG"
	infoLevel  = "INFO"
	warnLevel  = "WARN"
	errorLevel = "EERO"
)

var filePrefixFunc = func() string {
	abs, _ := filepath.Abs(".")
	return filepath.Dir(abs)
}()

func fileLine() string {
	pc, absPath, line, _ := runtime.Caller(2)
	caller := runtime.FuncForPC(pc)
	relPath := absPath[len(filePrefixFunc)+1:]
	simpleCaller := caller.Name()[len("ezmq."):]
	return fmt.Sprintf("%s:%v %s(): ", relPath, line, simpleCaller)
}
