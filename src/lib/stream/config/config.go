package config

import (
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
