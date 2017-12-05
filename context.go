package goutils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/rs/xid"
)

// ServerContext 日志上下文
type ServerContext struct {
	lock  *sync.Mutex
	buf   *bytes.Buffer
	uuid  string
	sTime time.Time
	tTime time.Time //临时统计用的时间
}

// NewContext 构造函数
func NewContext(msg string) *ServerContext {
	sc := new(ServerContext)
	sc.buf = bytes.NewBufferString(msg)
	sc.uuid = xid.New().String()
	sc.sTime = time.Now()
	sc.lock = new(sync.Mutex)
	//sc.tTime
	return sc
}

// SetUUID 设置上下文uuid，用于trace整个工作流
func (sc *ServerContext) SetUUID(uuid string) {
	sc.lock.Lock()
	if len(uuid) != 0 {
		sc.uuid = uuid
	}
	sc.lock.Unlock()
}

// GetUUID 获取当前上下文uuid
func (sc *ServerContext) GetUUID() string {
	return sc.uuid
}

// StartTimer 调用开始计时，用于统计程序耗时，和StopTimer配合使用
func (sc *ServerContext) StartTimer() {
	sc.tTime = time.Now()
}

// StopTimer 结束计时，和StartTimer配合使用
func (sc *ServerContext) StopTimer(key string) {
	duration := time.Now().Sub(sc.tTime)
	sc.lock.Lock()
	sc.buf.WriteString(fmt.Sprintf(" %s=%v", key, duration))
	sc.lock.Unlock()
}

// AddNotes 添加kv对到日志中
func (sc *ServerContext) AddNotes(key string, val interface{}) {
	sc.lock.Lock()
	sc.buf.WriteString(fmt.Sprintf(" %s=%v", key, val))
	sc.lock.Unlock()
}

// Flush flush所有AddNotes日志，通常工作流结束调用
func (sc *ServerContext) Flush() {
	duration := time.Now().Sub(sc.sTime)
	sc.lock.Lock()
	bufStr := sc.buf.String()
	sc.lock.Unlock()
	Log.Info(fmt.Sprintf("%s=%s cost=%v %s ", "Uuid", sc.uuid, duration, bufStr))
}

// Debug debug日志
func (sc *ServerContext) Debug(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.uuid, format)
	Log.Debug(s, args...)
}

// Info Info日志
func (sc *ServerContext) Info(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.uuid, format)
	Log.Info(s, args...)
}

// Notice Notice日志
func (sc *ServerContext) Notice(format string, args ...interface{}) {
	s := fmt.Sprintf("%s=%s %s", "Uuid", sc.uuid, format)
	Log.Notice(s, args...)
}

// Warning Warning日志
func (sc *ServerContext) Warning(format string, args ...interface{}) {
	_, fileName, lineNo := getRuntime(2)
	s := fmt.Sprintf("%s=%s %s=%s:%d %s", "Uuid", sc.uuid, "Runtime", fileName, lineNo, format)
	Log.Warning(s, args...)
}

// Error Error日志
func (sc *ServerContext) Error(format string, args ...interface{}) {
	_, fileName, lineNo := getRuntime(2)
	s := fmt.Sprintf("%s=%s %s=%s:%d %s", "Uuid", sc.uuid, "Runtime", fileName, lineNo, format)
	Log.Error(s, args...)
}

// Critical Critical日志
func (sc *ServerContext) Critical(format string, args ...interface{}) {
	_, fileName, lineNo := getRuntime(2)
	s := fmt.Sprintf("%s=%s %s=%s:%d %s", "Uuid", sc.uuid, "Runtime", fileName, lineNo, format)
	Log.Critical(s, args...)
}

func getRuntime(skip int) (function, filename string, lineno int) {
	function = "???"
	pc, filename, lineno, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}
	filename = filepath.Base(filename)
	return
}
