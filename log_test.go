package goutils

import (
	"testing"
)

func Test_LogInit(t *testing.T) {
	InitLog(nil)
	Log.Info("format")
}

func Test_LogDefault(t *testing.T) {
	df := new(defaultLogger)
	InitLog(df)
	Log.Info("format")
}
