package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig() Conf {
	var path string = cliArgs()

	var confData []byte

	if path != "" {
		confData = readYaml(path)
	} else {
		confData = []byte(defaultYaml)
	}

	conf := Conf{}

	err := yaml.Unmarshal(confData, &conf)
	if err != nil {
		panic(err)
	}

	log.Println(conf)

	for c, _ := range conf.Components {
		log.Printf("Component %#v, enabled: %#v\n", c, conf.Components[c])
	}

	for a, _ := range conf.Associations {
		for ac, _ := range conf.Associations[a] {
			log.Printf("Association %#v, item: %#v", a, conf.Associations[a][ac])
		}
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
