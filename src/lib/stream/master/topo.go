package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadConfig(configFileName string) (bolts *Bolts) {
	fileContent, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println("Cannot read the topology config file", err)
		return
	}
	cfg := new(Config)
	if err := json.Unmarshal(fileContent, cfg); err != nil {
		fmt.Println("Fail to parse the JSON topology file")
		return
	}
	return Config2Bolts(cfg)
}

func Config2Bolts(config *Config) *Bolts {
	oBolts := make([]Bolt, len(config.Bolts))
	for i, bolt := range config.Bolts {
		oBolts[i].ID = bolt.ID
		oBolts[i].Name = bolt.Name
		oBolts[i].Pred = bolt.Pred
	}
	return &Bolts{Bolts: oBolts}
}

/*
	I change the structure of Config, with independent bolt struct to allow
	mode convenience
*/
type Bolts struct {
	Bolts []Bolt
}

type Bolt struct {
	ID   int
	Name string
	Pred []int
}

type Config struct {
	Bolts []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Pred []int  `json:"pred"`
	} `json:"bolts"`
}
