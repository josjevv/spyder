package config

type Conf struct {
	MongoHost    string
	MongoDb      string
	Components   map[string]bool
	Associations map[string][]string
}

func (c Conf) HasComponent(component string) bool {
	_, present := c.Components[component]
	return present
}

func (c Conf) HasAssociation(association string) bool {
	if _, present := c.Associations["all"]; present {
		return true
	}
	_, present := c.Associations[association]
	return present
}
