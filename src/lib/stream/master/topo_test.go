package master

import "testing"

func TestReadConfig(t *testing.T) {
	filename := "../../../../examples/streamProcessing/exclamation/topo.json"
	config := ReadConfig(filename)
	println(config.Bolts[0].Name)
}
