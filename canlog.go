// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package canlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// These flags define what level message to write
const (
	LZero = iota
	LDebug
	LInfo
	LWarn
	LError
	LFatal
)

var levelPrefix = []string{"[     ]", "[Debug]", "[ Info]", "[ Warn]", "[Error]", "[Fatal]"}

var logLevel = 0
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// SetWriter sets the destination on which log data will be written.
// The prefix appears at the beginning of each line followed a space.
func SetWriter(rw io.Writer, prefix string) {
	if !strings.HasSuffix(prefix, " ") {
		prefix = prefix + " "
	}
	// todo 性能和格式
	logger = log.New(rw, prefix, log.LstdFlags|log.Lshortfile)
}

// GetLogger returns the log.Logger
func GetLogger() *log.Logger {
	return logger
}

func canLine(level int, v ...interface{}) {
	if level >= logLevel {
		_ = logger.Output(3, levelPrefix[level]+" "+fmt.Sprintln(v...))
	}
}

// CanOutput writes the output for a logging line.
// The str contains the text to print after prefix and level-prefix.
// callDepth is used to recover the PC adn is provided for generality.
func CanOutput(callDepth int, level int, str string) {
	_ = logger.Output(callDepth, levelPrefix[level]+" "+str)
}

// CanDebug call CanOutput with LDebug
func CanDebug(v ...interface{}) {
	canLine(LDebug, v...)
}

// CanInfo call CanOutput with LInfo
func CanInfo(v ...interface{}) {
	canLine(LInfo, v...)
}

// CanWarn call CanOutput with LWarn
func CanWarn(v ...interface{}) {
	canLine(LWarn, v...)
}

// CanError call CanOutput with LError
func CanError(v ...interface{}) {
	canLine(LError, v...)
}

// CanFatal will panic
func CanFatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
