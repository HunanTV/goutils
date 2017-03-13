package goutils

import "testing"

func Test_NewContext(t *testing.T) {
	context := NewContext("test")
	if context == nil {
		t.Fail()
	}
}

func Test_UUID(t *testing.T) {
	context := NewContext("test")
	if context == nil {
		t.Fail()
	}
	context.SetUUID("test")
	if "test" != context.GetUUID() {
		t.Fail()
	}
	context.Notice("format")
}

func Test_AddNotesAndFlush(t *testing.T) {
	context := NewContext("test")
	if context == nil {
		t.Fail()
	}
	context.AddNotes("test", 1)
	context.Flush()
}

func Test_LogLevel(t *testing.T) {
	InitLog(nil)
	context := NewContext("test")
	context.Debug("you will not see me")
	context.Info("you will see me")

}

func Test_Log(t *testing.T) {
	context := NewContext("test")
	context.Debug("debug")
	context.Info("info")
	context.Warning("warnging")
	context.Error("error")
	context.Critical("critical")
	context.Notice("Notice")
}
