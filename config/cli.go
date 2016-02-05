package config

import (
	"flag"
	"github.com/bulletind/spyder/log"
)

func cliArgs() string {
	var config = flag.String(
		"config",
		"",
		"Config [.yml format] file to load the configurations from",
	)

	flag.Parse()

	if *config == "" {
		log.Info("No config file supplied. Using defauls.")
	}

	return *config
}
