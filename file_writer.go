// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package canlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type fileWriter struct {
	locker      sync.RWMutex
	fileName    string
	file        *os.File
	openDate    string
	endTime     time.Time
	rotateChan  chan string
	errorLogger *log.Logger
}

func NewFileWriter(fileName string) io.Writer {
	return initFileWriter(new(fileWriter), fileName)
}

func initFileWriter(fw *fileWriter, fileName string) *fileWriter {
	fw.errorLogger = log.New(os.Stderr, "file_writer", log.Llongfile|log.LstdFlags)
	fw.rotateChan = make(chan string, 1)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fw.errorLogger.Println(err)
		// em???
	}
	var openTime time.Time
	if fw.endTime.IsZero() {
		openTime = time.Now()
	} else {
		// 文件按天rotate 需要加1天
		openTime = fw.endTime.AddDate(0, 0, 1)
	}
	openDate := openTime.Format("2006-01-02")
	endTime := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 24, 0, 0, 0, openTime.Location())

	fw.fileName = fileName
	fw.openDate = openDate
	fw.endTime = endTime
	fw.file = file
	go fw.watchRotate()
	return fw
}

func (fw *fileWriter) Write(p []byte) (n int, err error) {
	fw.locker.Lock()
	n, err = fw.file.Write(p)
	fw.locker.Unlock()
	return n, err
}

func (fw *fileWriter) watchRotate() {
	time.AfterFunc(fw.endTime.Sub(time.Now()), func() {
		fw.rotateChan <- fw.fileName + "-" + fw.openDate
	})
	go func() {
		fileName := <-fw.rotateChan
		fw.locker.Lock()
		var err error
		err = fw.file.Close()
		if err != nil {
			// em???
			fw.errorLogger.Println(err)
		}
		for {
			fileInfo, _ := os.Lstat(fileName)
			if fileInfo != nil {
				fileName = fmt.Sprintf("%s-%d", fileName, time.Now().UnixNano())
			} else {
				break
			}
		}
		err = os.Rename(fw.fileName, fileName)
		if err != nil {
			// em???
			fw.errorLogger.Println(err)
		}
		initFileWriter(fw, fw.fileName)
		fw.locker.Unlock()
	}()
}
