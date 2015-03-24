package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig(conf interface{}, defaultYaml string) {
	var path string = cliArgs()
	err := yaml.Unmarshal([]byte(defaultYaml), conf)
	if err != nil {
		panic(err)
	}

	if path != "" {
		confData := readYaml(path)
		err := yaml.Unmarshal(confData, conf)
		if err != nil {
			panic(err)
		}
	}
}

func readYaml(path string) []byte {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return yamlFile
}
