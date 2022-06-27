package logger

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kataras/golog"
)

var (
	_default   = golog.Default.SetLevel("info")
	filePrefix = func() string {
		abs, _ := filepath.Abs(".")
		return filepath.Dir(abs)
	}()
)

func init() {
	_default.SetTimeFormat("2006/01/02 15:04:05")
}

func LevelString(level golog.Level) string {
	switch level {
	case golog.DisableLevel:
		return "disable"
	case golog.FatalLevel:
		return "fatal"
	case golog.ErrorLevel:
		return "error"
	case golog.WarnLevel:
		return "warn"
	case golog.InfoLevel:
		return "info"
	case golog.DebugLevel:
		return "debug"
	}
	return ""
}

func checkLevel(lv string) error {
	switch lv {
	case "disable", "fatal", "error", "warn", "info", "debug":
		return nil
	default:
		return errors.New("doesn't exist Log Level: " + lv)
	}
}

func Level() golog.Level {
	return _default.Level
}

// Level in "disable", "fatal", "error", "warn", "info", "debug"
func SetLevel(lv string) {
	err := checkLevel(lv)
	if err != nil {
		Info("using default log level INFO")
		_default.SetLevel("info")
		return
	}
	_default.SetLevel(lv)
	Info("using ", strings.ToUpper(lv), " log level")
}

func fileLine() string {
	pc, absPath, line, _ := runtime.Caller(2)
	caller := runtime.FuncForPC(pc)
	relPath := absPath[len(filePrefix)+1:]
	simpleCaller := caller.Name()[len("ezmq."):]
	return fmt.Sprintf("%s:%v %s(): ", relPath, line, simpleCaller)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(v ...interface{}) {
	nv := append([]interface{}{fileLine()}, v...)
	_default.Fatal(nv...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	_default.Fatalf(fileLine()+format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(v ...interface{}) {
	nv := append([]interface{}{fileLine()}, v...)
	_default.Error(nv...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) {
	_default.Errorf(fileLine()+format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(v ...interface{}) {
	nv := append([]interface{}{fileLine()}, v...)
	_default.Warn(nv...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	_default.Warnf(fileLine()+format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(v ...interface{}) {
	nv := append([]interface{}{fileLine()}, v...)
	_default.Info(nv...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	_default.Infof(fileLine()+format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(v ...interface{}) {
	nv := append([]interface{}{fileLine()}, v...)
	_default.Debug(nv...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	_default.Debugf(fileLine()+format, args...)
}
