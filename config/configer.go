package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig() Conf {
	var path string = cliArgs()
	conf := Conf{}
	yaml.Unmarshal([]byte(defaultYaml), &conf)

	if path != "" {
		confData := readYaml(path)
		yaml.Unmarshal(confData, &conf)
	}

	return conf
}

func readYaml(path string) []byte {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return yamlFile
}
