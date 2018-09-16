package CS425_LOG

import (
	"fmt"
	"log"
	"os"
)

const (
	NAME = "CS425-MP1-LOGFILE-VM-%d"
)

type LogMessage struct {
	vmIndex int32
	logMode int32
	file    *os.File
	err     error
	logger  *log.Logger
}

func (T *LogMessage) Print(str string) {
	T.logger.Println(str)
}

func (T *LogMessage) Close() {
	T.file.Close()
	log.Printf("Log File for server %d has closed.\n", T.vmIndex)
}

func (T *LogMessage) Init(index, mode int32) {
	T.vmIndex = index
	T.logMode = mode
	FileName := fmt.Sprintf(NAME, T.vmIndex)
	T.file, T.err = os.Create(FileName)
	if T.err != nil {
		fmt.Println("Unable to create the file.")
	}
	//log using citing from website https://www.jianshu.com/p/8e6dd09e86f7
	T.logger = log.New(T.file, "", log.LstdFlags)
}
