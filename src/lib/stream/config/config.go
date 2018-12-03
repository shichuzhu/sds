package config

import (
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/cfg"
	"fa18cs425mp/src/shared/sdfs2fd"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
)

var RootPathp = flag.String("stream", cfg.Cfg.StreamDir, "Directory to store the Crane streaming data, RELATIVE to GOPATH/src")
var RootPath string
var RootPathRelativeToGoPath string // used when 'go build'

func InitialCrane() {
	sdfs2fd.Fd2Crane = make(chan int, 2)
	gopath := func() string {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}
		return gopath
	}()
	RootPathRelativeToGoPath = *RootPathp

	RootPath = gopath + "/src/" + RootPathRelativeToGoPath

	println(RootPath, RootPathRelativeToGoPath)
	_ = os.RemoveAll(RootPath)
	log.Println("Creating directory: ", RootPath)
	err := os.MkdirAll(RootPath, os.ModePerm)
	if err != nil {
		log.Println("err MkdirAll: ", err)
	}
}

func SetupLoadFile(jobName string) (err error) {
	err = os.RemoveAll(RootPath + jobName)
	err = os.MkdirAll(RootPath+jobName, os.ModePerm)
	dirpath := RootPath + jobName + "/"
	//log.Println("local job directory: ", dirpath)
	_ = os.Mkdir(dirpath+"plugin/", os.ModePerm)
	_ = os.Mkdir(dirpath+"src/", os.ModePerm)

	var cmd string

	// TODO: load file from sdfs
	_ = utils.RunShellString(fmt.Sprintf("sds sdfs get %s.zip %ssrc/%s.zip", jobName, dirpath, jobName))
	//_ = utils.RunShellString("zip -rj data/mp4/exclamation/src/exclamation.zip examples/streamProcessing/exclamation")

	cmd = "unzip -od " + dirpath + "src " + dirpath + "src/" + jobName + ".zip"
	err = utils.RunShellString(cmd)
	return err
}

func DirPath(jobName string) string {
	return RootPath + jobName + "/"
}
