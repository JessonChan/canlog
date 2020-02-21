package canlog

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type fileWriter struct {
	locker   sync.RWMutex
	fileName string
	file     *os.File
	openTime time.Time
	openDate string
	endTime  time.Time
}

func NewFileWriter(fileName string) io.Writer {
	return newFileWriter(new(fileWriter), fileName)
}

var rotateChan = make(chan string, 1)
var errLogger = log.New(os.Stderr, "file_writer", log.Llongfile|log.LstdFlags)

func newFileWriter(fw *fileWriter, fileName string) *fileWriter {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		errLogger.Println(err)
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
		rotateChan <- fw.fileName
	})
	go func() {
		fileName := <-rotateChan
		fw.locker.Lock()
		var err error
		err = fw.file.Close()
		if err != nil {
			// em???
			errLogger.Println(err)
		}
		// todo 如果 fw.fileName+"-"+fw.openDate 已经存在
		err = os.Rename(fw.fileName, fw.fileName+"-"+fw.openDate+fileName)
		if err != nil {
			// em???
			errLogger.Println(err)
		}
		newFileWriter(fw, fw.fileName)
		fw.locker.Unlock()
	}()
}
