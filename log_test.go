package goutils

import (
	"testing"
)

func Test_LogInit(t *testing.T) {
	InitLog("", "INFO")
	Log.Info("format")
}
