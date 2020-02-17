package canlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	LZero = iota
	LDebug
	LInfo
	LWarn
	LError
	LFatal
)

var logPrefix = []string{"[     ]", "[Debug]", "[ Info]", "[ Warn]", "[Error]", "[Fatal]"}

var logLevel = 0
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func SetWriter(rw io.Writer, prefix string) {
	if !strings.HasSuffix(prefix, " ") {
		prefix = prefix + " "
	}
	logger = log.New(rw, prefix, log.LstdFlags|log.Lshortfile)
}
func GetLogger() *log.Logger {
	return logger
}

func canLine(level int, v ...interface{}) {
	if level >= logLevel {
		_ = logger.Output(3, logPrefix[level]+" "+fmt.Sprintln(v...))
	}
}

func CanOutput(callDepth int, str string) {
	_ = logger.Output(callDepth, str)
}

func CanDebug(v ...interface{}) {
	canLine(LDebug, v...)
}

func CanInfo(v ...interface{}) {
	canLine(LInfo, v...)
}
func CanWarn(v ...interface{}) {
	canLine(LWarn, v...)
}
func CanError(v ...interface{}) {
	canLine(LError, v...)
}
func CanFatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
