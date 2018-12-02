package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadConfig(configFileName string) (Cfg *Config) {
	fileContent, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println("Cannot read the topology config file")
		return
	}
	if err := json.Unmarshal(fileContent, Cfg); err != nil {
		fmt.Println("Fail to parse the JSON topology file")
		return
	}
	return
}

type Config struct {
	/*Bolts []struct {
		ID   int           `json:"id"`
		Name string        `json:"name"`
		Pred []interface{} `json:"pred"`
	} `json:"bolts"`*/
	Bolts []Bolt
}

type Bolt struct {
	ID   int
	Name string
	pred []interface{}
}
