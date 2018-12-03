package config

import (
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/cfg"
	"flag"
	"os"
)

var RootPathp = flag.String("stream", cfg.Cfg.StreamDir, "Directory to store the Crane streaming data")
var RootPath string

func InitialCrane() {
	RootPath = *RootPathp
	_ = os.RemoveAll(RootPath)
	_ = os.Mkdir(RootPath, os.ModePerm)
}

func SetupLoadFile(jobName string) error {
	_ = os.RemoveAll(RootPath + jobName)
	_ = os.Mkdir(RootPath+jobName, os.ModePerm)
	dirpath := RootPath + jobName + "/"
	//log.Println("local job directory: ", dirpath)
	_ = os.Mkdir(dirpath+"plugin/", os.ModePerm)
	_ = os.Mkdir(dirpath+"src/", os.ModePerm)

	var cmd string

	// TODO: load file from sdfs
	//_ = utils.RunShellString(fmt.Sprintf("sds sdfs get %s.zip %ssrc/%s.zip", jobName, dirpath, jobName))
	_ = utils.RunShellString("zip -rj data/mp4/exclamation/src/exclamation.zip examples/streamProcessing/exclamation")

	cmd = "unzip -d " + dirpath + "src " + dirpath + "src/" + jobName + ".zip"
	err := utils.RunShellString(cmd)
	return err
}

func DirPath(jobName string) string {
	return RootPath + jobName + "/"
}
