package canlog

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func Test_fileWriter_Write(t *testing.T) {
	fileName := "/tmp/canlog.txt"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	err = os.Rename(fileName, fileName+"_bak")
	if err != nil {
		t.Error(err)
	}
	_, err = file.Write([]byte("ok" + time.Now().Format("15:04:05")))
	if err != nil {
		t.Error(err)
	}
	err = file.Close()
	if err != nil {
		t.Error(err)
	}
}

func Test_file(t *testing.T) {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}()
	time.AfterFunc(10*time.Second, func() {
		fmt.Println("hello")
	})
	time.Sleep(15 * time.Second)
}

func Test_Write(t *testing.T) {
	fileName := "/tmp/canlog.txt"
	fw := newFileWriter(new(fileWriter), fileName)
	l := log.New(fw, "TEST", log.Lshortfile|log.LstdFlags)
	go func() {
		for i := 0; ; i++ {
			l.Println("test", i)
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(500)))
		}
	}()
	for i := 0; ; i++ {
		time.Sleep(time.Second * 3)
		rotateChan <- fmt.Sprintf("%d", i)
	}
}
