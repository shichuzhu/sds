package worker

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/pb"
	"log"
	"plugin"
)

func CompilePlugin(cfg *pb.TaskCfg) *plugin.Plugin {
	var cmd string
	dirpath := config.RootPath + cfg.JobName + "/"

	filepath := dirpath + "plugin/" + cfg.JobName + ".so"
	cmd = "go build -buildmode=plugin -o " + filepath + " fa18cs425mp/" + dirpath + "src/"
	_ = utils.RunShellString(cmd)
	plug, err := plugin.Open(filepath)
	if err != nil {
		log.Println("error: ", err)
		return nil
	} else {
		return plug
	}
}

func SpawnBoltTask(cfg *pb.TaskCfg, plug *plugin.Plugin) shared.BoltABC {
	sym, err := plug.Lookup("New" + cfg.Bolt.Name)
	if err != nil {
		log.Println(err)
		return nil
	}
	var bolt shared.BoltABC
	bolt = sym.(func() shared.BoltABC)()
	return bolt
}

func SpawnSinkTask(cfg *pb.TaskCfg, plug *plugin.Plugin) shared.SinkABC {
	sym, err := plug.Lookup("New" + cfg.Bolt.Name)
	if err != nil {
		log.Println(err)
		return nil
	}
	var bolt shared.SinkABC
	bolt = sym.(func() shared.SinkABC)()
	return bolt
}

func SpawnSpoutTask(cfg *pb.TaskCfg, plug *plugin.Plugin) shared.SpoutABC {
	sym, err := plug.Lookup("New" + cfg.Bolt.Name)
	if err != nil {
		log.Println(err)
		return nil
	}
	var spout shared.SpoutABC
	spout = sym.(func() shared.SpoutABC)()
	return spout
}
