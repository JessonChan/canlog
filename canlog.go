// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package canlog

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logger = newColorConsole()

func newColorConsole() *CanLogger {
	cl := NewCanLogger(os.Stdout, "")
	cl.isColor = true
	return cl
}

// SetWriter sets the destination on which log data will be written.
// The prefix appears at the beginning of each line followed a space.
func SetWriter(rw io.Writer, prefix string) {
	if rw != os.Stdout {
		logger.isColor = false
	}
	logger.SetOutput(rw)
	logger.SetPrefix(formatPrefix(prefix))
}

// SetPrefix sets the output prefix(followed by a space) for the logger.
func SetPrefix(prefix string) {
	logger.SetPrefix(formatPrefix(prefix))
}

func SetLogLevel(level int) {
	logger.SetLogLevel(level)
}

// GetLogger returns the log.Logger
func GetLogger() *log.Logger {
	return logger.Logger
}

// CanOutput writes the output for a logging line.
// The str contains the text to print after prefix and level-prefix.
// callDepth is used to recover the PC adn is provided for generality.
func CanOutput(callDepth int, level int, str string) {
	logger.CanOutput(callDepth, level, str)
}

// CanDebug call CanOutput with LDebug
func CanDebug(v ...interface{}) {
	logger.canLine(LDebug, v...)
}

// CanInfo call CanOutput with LInfo
func CanInfo(v ...interface{}) {
	logger.canLine(LInfo, v...)
}

// CanWarn call CanOutput with LWarn
func CanWarn(v ...interface{}) {
	logger.canLine(LWarn, v...)
}

// CanError call CanOutput with LError
func CanError(v ...interface{}) {
	logger.canLine(LError, v...)
}

// CanFatal will panic
func CanFatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
