package immap

import (
	"context"
)

type baseMap struct {
	context.Context
	addChan   chan *mapPack
	checkChan chan *mapPack
	delChan   chan *mapPack
}

type contextMap baseMap

type ImutMap baseMap

type retPack struct {
	value  interface{}
	mapRef IntfMap
}

type mapPack struct {
	key, value interface{}
	ret        chan retPack
}

// contextMapper implements the lockless map interface.
type contextMapper interface {
	Exists(interface{}) (interface{}, bool)
	Add(interface{}, interface{}) error
	Delete(interface{})
}

type ImutMapper interface {
	Exists(interface{}) (interface{}, bool, IntfMap)
	Add(interface{}, interface{}) (IntfMap, error)
	Delete(interface{})
}
