package spyder

type FlyCallback func(fly Fly)

var FlyRegistry = map[string]FlyCallback{}
