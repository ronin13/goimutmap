package immap

import (
	"context"
)

type ImMap struct {
	context.Context
	addChan   chan *ImPack
	checkChan chan *ImPack
	done      chan struct{}
	stopMap   context.CancelFunc
}

type ImutMap struct {
	context.Context
	addChan   chan *ImPack
	checkChan chan *ImPack
	done      chan struct{}
	stopMap   context.CancelFunc
}

type retMap struct {
	value  interface{}
	mapRef IntfMap
}

type ImPack struct {
	key, value interface{}
	ret        chan retMap
}

// ImMapper implements the lockless map interface for use by crawler.
type ImMapper interface {
	Exists(interface{}) (interface{}, bool)
	Add(interface{}, interface{}) error
	Stop()
	// Delete
}

type ImutMapper interface {
	Exists(interface{}) (interface{}, bool, IntfMap)
	Add(interface{}, interface{}) (IntfMap, error)
	Stop()
	// Delete
}
