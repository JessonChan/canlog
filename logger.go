// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package canlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
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
var currentPath = func() string {
	dir, err := os.Getwd()
	if err == nil {
		if len(dir) > 0 {
			if dir[len(dir)-1] != '/' {
				dir = dir + "/"
			}
			return dir
		}
	}
	return ""
}()
var currentPathLen = len(currentPath)

var logLevel = 0

type CanLogger struct {
	isColor bool
	*log.Logger
}

func formatPrefix(p string) string {
	if !strings.HasSuffix(p, " ") {
		p = p + " "
	}
	return p
}

func NewCanLogger(rw io.Writer, prefix string) *CanLogger {
	return &CanLogger{isColor: false, Logger: log.New(rw, formatPrefix(prefix), log.LstdFlags)}
}

var colorBrush = []func(format string, a ...interface{}) string{color.RedString, color.GreenString, color.YellowString, color.BlueString, color.MagentaString, color.CyanString, color.WhiteString}

func (cl *CanLogger) canLine(level int, v ...interface{}) {
	if cl.isColor {
		cs := ""
		for i := 0; i < len(v); i++ {
			cs = cs + colorBrush[i%len(colorBrush)](fmt.Sprint(v[i]))
		}
		v = []interface{}{cs}
	}
	if level >= logLevel {
		cl.CanOutput(3, level, fmt.Sprintln(v...))
	}
}

// CanOutput writes the output for a logging line.
// The str contains the text to print after prefix and level-prefix.
// callDepth is used to recover the PC adn is provided for generality.
func (cl *CanLogger) CanOutput(callDepth int, level int, str string) {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		file = "???"
		line = 0
	} else {
		short := file
		cd := callDepth
		if len(currentPath) <= len(file) {
			if file[currentPathLen-1] == '/' {
				cd = -1
				short = file[currentPathLen:]
			}
		}
		if cd >= 0 {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					cd--
					if cd == 0 {
						short = file[i+1:]
						break
					}
				}
			}
		}
		file = short
	}
	if cl.isColor {
		_ = cl.Output(callDepth, color.GreenString(levelPrefix[level])+" "+color.BlueString(file)+":"+color.CyanString(fmt.Sprintf("%d ", line))+str)
	} else {
		_ = cl.Output(callDepth, levelPrefix[level]+" "+file+":"+fmt.Sprintf("%d ", line)+str)
	}
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
