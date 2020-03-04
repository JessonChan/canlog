// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package canlog

import (
	"fmt"
	"log"
)

type CanLogger struct {
	*log.Logger
}

func (cl *CanLogger) canLine(level int, v ...interface{}) {
	if level >= logLevel {
		_ = cl.Output(3, levelPrefix[level]+" "+fmt.Sprintln(v...))
	}
}

// CanOutput writes the output for a logging line.
// The str contains the text to print after prefix and level-prefix.
// callDepth is used to recover the PC adn is provided for generality.
func (cl *CanLogger) CanOutput(callDepth int, level int, str string) {
	_ = cl.Output(callDepth, levelPrefix[level]+" "+str)
}

// CanDebug call CanOutput with LDebug
func (cl *CanLogger) CanDebug(v ...interface{}) {
	cl.canLine(LDebug, v...)
}

// CanInfo call CanOutput with LInfo
func (cl *CanLogger) CanInfo(v ...interface{}) {
	cl.canLine(LInfo, v...)
}

// CanWarn call CanOutput with LWarn
func (cl *CanLogger) CanWarn(v ...interface{}) {
	cl.canLine(LWarn, v...)
}

// CanError call CanOutput with LError
func (cl *CanLogger) CanError(v ...interface{}) {
	cl.canLine(LError, v...)
}

// CanFatal will panic
func (cl *CanLogger) CanFatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
