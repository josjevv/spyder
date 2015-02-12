package spyder

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

func readConfig() Config {
  filename, _ := filepath.Abs("./config.yml")
  yamlFile, err := ioutil.ReadFile(filename)

  if err != nil {
    panic(err)
  }

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
  Components map[string]bool
  Associations map[string][]string
}
