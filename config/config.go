package config

type Conf struct {
	MongoHost     string
	MongoDb       string
	Components    map[string]bool
	Associations  map[string][]string
	Notifications map[string]bool
	History       map[string]bool
}
