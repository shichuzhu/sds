package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	NAME = "vm-%d.log"
)

type LogMessage struct {
	file   *os.File
	err    error
	logger *log.Logger
}

func (T *LogMessage) Print(str string) {
	T.logger.Println(str)
}

func (T *LogMessage) Close() {
	T.file.Close()
	log.Println("Log File has closed.")
}

func (T *LogMessage) InitLogger(index int, path string) {
	FileName := fmt.Sprintf(path+"/"+NAME, index)
	T.file, T.err = os.Create(FileName)
	if T.err != nil {
		fmt.Println("Unable to create the file.")
	}
	//log using citing from website https://www.jianshu.com/p/8e6dd09e86f7
	//T.logger = log.New(T.file, "", log.LstdFlags)
	log.SetOutput(T.file)
	log.SetPrefix(fmt.Sprintf("N%02d: ", index))
}
