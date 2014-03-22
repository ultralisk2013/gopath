package util

import (
	"fmt"
	"time"
)

const (
	SHOW_INFO     = true
	SHOW_DEBUG    = true
	SHOW_WARNNING = true
	SHOW_ERROR    = true
	TIME_FMT      = "[06-01-02 15:04:05.000]"
)

func TimeString(timeValue int64) string {
	var t time.Time
	if timeValue != 0 {
		t = time.Unix(timeValue, 0)
	} else {
		t = time.Now()
	}
	return t.Format(TIME_FMT)
}

func ShowInfo(a ...interface{}) {
	if SHOW_INFO {
		//fmt.Println("\033[30m[INFO]\033[0m", TimeString(0), a)
		fmt.Println("[INFO]", TimeString(0), a)
	}
}
func ShowDebug(a ...interface{}) {
	if SHOW_DEBUG {
		fmt.Println("[DEBUG]", TimeString(0), a)
	}
}

func ShowWarnning(a ...interface{}) {
	if SHOW_WARNNING {
		fmt.Println("[WARN]", TimeString(0), a)
	}
}

func ShowError(a ...interface{}) {
	if SHOW_ERROR {
		fmt.Println("[ERROR]", TimeString(0), a)
	}
}
func CaoSiShowDebug(a ...interface{}) {
	if SHOW_DEBUG {
		fmt.Println("[CAOSI_DEBUG]", TimeString(0), a)
	}
}
