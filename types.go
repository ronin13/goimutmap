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

type ImPack struct {
	key, value interface{}
	ret        chan interface{}
}

// ImMapper implements the lockless map interface for use by crawler.
type ImMapper interface {
	Exists(interface{}) (interface{}, bool)
	Add(interface{}, interface{}) error
	Stop()
	// Delete
}
