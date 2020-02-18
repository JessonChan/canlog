package canlog

import (
	"io"
	"os"
	"time"
)

type fileWriter struct {
	fileName string
	file     *os.File
	openTime time.Time
	openDate string
	endTime  time.Time
}

func NewFileWriter(fileName string) io.Writer {
	return newFileWriter(fileName)
}

func newFileWriter(fileName string) *fileWriter {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		// em???
	}
	openTime := time.Now()
	openDate := openTime.Format("2006-01-02")
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", openDate+" 23:59:59", openTime.Location())
	fw := &fileWriter{fileName: fileName, file: file, openTime: openTime, openDate: openDate, endTime: endTime}
	go fw.rotate()
	return fw
}

func (fw *fileWriter) Write(p []byte) (n int, err error) {
	return fw.file.Write(p)
}

func (fw *fileWriter) rotate() {
	time.AfterFunc(fw.endTime.Sub(fw.openTime), func() {
		err := os.Rename(fw.fileName, fw.fileName+"-"+fw.openDate)
		if err != nil {
			// em???
		}
		nfw := newFileWriter(fw.fileName)
		err = fw.file.Close()
		if err != nil {
			// em???
		}
		fw.file = nfw.file
		fw.openDate = nfw.openDate
		fw.openTime = nfw.openTime
		fw.endTime = nfw.endTime
		fw.rotate()
	})
}
