package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Connstring   string
	Components   map[string]bool
	Associations map[string][]string
}

func ReadConfig() Config {
	var path string = CliArgs()
	var config Config

	if path != "" {
		config = readYaml(path)
	} else {
		config = GetDefaultConfig()
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

func readYaml(path string) Config {
	var config Config

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
		panic(err)
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
