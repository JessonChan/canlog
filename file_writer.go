// Copyright 2020 JessonChan.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package canlog

import (
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
	openTime    time.Time
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
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fw.errorLogger.Println(err)
		// em???
	}
	openTime := time.Now()
	openDate := openTime.Format("2006-01-02")
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", openDate+" 23:59:59", openTime.Location())
	fw.fileName = fileName
	fw.openTime = openTime
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
	time.AfterFunc(fw.endTime.Sub(fw.openTime), func() {
		fw.rotateChan <- fw.fileName + "_" + fw.openDate
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
		// todo 如果 fw.fileName+"-"+fw.openDate 已经存在
		err = os.Rename(fw.fileName, fileName)
		if err != nil {
			// em???
			fw.errorLogger.Println(err)
		}
		initFileWriter(fw, fw.fileName)
		fw.locker.Unlock()
	}()
}
