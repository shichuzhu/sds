package worker

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/pb"
	"log"
	"os"
	"plugin"
)

func SetupDirectories(cfg *pb.TaskCfg) {
	_ = os.RemoveAll(config.RootPath + cfg.JobName)
	_ = os.Mkdir(config.RootPath+cfg.JobName, os.ModePerm)
	dirpath := config.RootPath + cfg.JobName + "/"
	_ = os.Mkdir(dirpath+"plugin/", os.ModePerm)
	_ = os.Mkdir(dirpath+"src/", os.ModePerm)
}

func CompilePlugin(cfg *pb.TaskCfg) *plugin.Plugin {
	var cmd string
	dirpath := config.RootPath + cfg.JobName + "/"

	cmd = "sds sdfs get " + cfg.JobName + ".zip " + dirpath + "src/" + cfg.JobName + ".zip"
	_ = utils.RunShellString(cmd)
	cmd = "unzip -d " + dirpath + "src " + dirpath + "src/" + cfg.JobName + ".zip"
	_ = utils.RunShellString(cmd)
	filepath := dirpath + "plugin/" + cfg.JobName + ".so"
	cmd = "go build -buildmode=plugin -o " + filepath + " fa18cs425mp/" + dirpath + "src/" + cfg.JobName
	_ = utils.RunShellString(cmd)
	plug, err := plugin.Open(filepath)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return plug
	}
}

func SpawnTask(cfg *pb.TaskCfg, plug *plugin.Plugin) shared.SpoutABC {
	sym, err := plug.Lookup("NewSpout")
	if err != nil {
		log.Println(err)
		return nil
	}
	var spout shared.SpoutABC
	spout = sym.(func() shared.SpoutABC)()
	return spout
}
