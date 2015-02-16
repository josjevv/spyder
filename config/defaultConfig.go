package config

//TODO - could or should we use just a string or json here?
func GetDefaultConfig() Config {
	return Config{
		Connstring: "localhost",
		Components: map[string]bool{
			"history":       true,
			"notifications": false},
		Associations: map[string][]string{
			"incidents": []string{"history", "notifications"}}}
}

// yaml file content
// connstring: "localhost"

// components:
//     history: true
//     notifications: false

// associations:
//     incidents: [history, notifications]
