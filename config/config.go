package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

func ReadConfig(path string) Config {
	var yamlFile []byte
	var err error

	if path != "" {
		path, _ = filepath.Abs("./default.yml")
	}
	if path != "" {
		yamlFile, err = ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
	} //else {
	//yamlFile := defaultConfig //Its a string not a location to a file.
	//}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	for c, _ := range config.Components {
		fmt.Printf("Component %#v, enabled: %#v\n", c, config.Components[c])
	}
	for a, _ := range config.Associations {
		for ac, _ := range config.Associations[a] {
			fmt.Printf("Association %#v, item: %#v\n", a, config.Associations[a][ac])
		}
	}
	return config
}

type Config struct {
	Components   map[string]bool
	Associations map[string][]string
}
