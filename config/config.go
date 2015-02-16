package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Config struct {
	Connstring   string
	Components   map[string]bool
	Associations map[string][]string
}

func ReadConfig() Config {
	var path string = CliArgs()
	var yamlFile []byte
	var err error

	if path == "" {
		path, _ = filepath.Abs("./config/default.yml")
	}
	if path != "" {
		yamlFile, err = ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	} //else {
	//yamlFile := defaultConfig //Its a string not a location to a file.
	//}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for c, _ := range config.Components {
		log.Printf("Component %#v, enabled: %#v\n", c, config.Components[c])
	}
	for a, _ := range config.Associations {
		for ac, _ := range config.Associations[a] {
			log.Printf("Association %#v, item: %#v\n", a, config.Associations[a][ac])
		}
	}
	return config
}

func CliArgs() string {
	var config = flag.String(
		"config",
		"",
		"Config [.yml format] file to load the configurations from",
	)

	flag.Parse()

	if *config == "" {
		log.Println("No config file supplied. Using defauls.")
	}

	return *config
}
