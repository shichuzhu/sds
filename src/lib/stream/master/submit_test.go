package master

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/pb"
	"flag"
	"testing"
)

func TestSubmitJob(t *testing.T) {
	flag.Parse()
	config.InitialCrane()
	_, _ = SubmitJob(&pb.TopoConfig{JobName: "exclamation"})
}
