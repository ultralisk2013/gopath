package util

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

func RuntimeInfo() (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(2)
}

// http发送简单文本
func HttpSendSimplePage(w *http.ResponseWriter, code int, content string) {
	if w != nil {
		(*w).WriteHeader(code)
		(*w).Write([]byte(content))
	}
}

// 发送空返回
func HttpSend200Empty(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 200, "empty")
}

// 发送404页面
func HttpSend404NotFound(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 404, "NOT FOUND")
}

func Error(fmtStr string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(fmtStr, a))
}
