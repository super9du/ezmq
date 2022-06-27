package logger

import (
	"testing"
)

func TestSetGlobalLevel(t *testing.T) {
	SetLevel("warn")
	Info("test")
	Warn("test")
}
