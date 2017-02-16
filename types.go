package immap

import (
	"context"
)

type OperType int

const (
	ADD_KEY OperType = iota
	CHECK_KEY
	DEL_KEY
)

type IntfMap map[interface{}]interface{}

type baseMap struct {
	context.Context
	cChan chan *mapPack
}

type contextMap baseMap

type ImutMap baseMap

type retPack struct {
	value  interface{}
	mapRef IntfMap
}

type mapPack struct {
	op         OperType
	key, value interface{}
	ret        chan retPack
}

// contextMapper implements the lockless map interface.
type ContextMapper interface {
	Exists(interface{}) (interface{}, bool)
	Add(interface{}, interface{}) interface{}
	Delete(interface{})
}

type ImutMapper interface {
	Exists(interface{}) (interface{}, bool, IntfMap)
	Add(interface{}, interface{}) IntfMap
	Delete(interface{})
}
