package immap

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

// NewImMapper returns a new instance of implementing ImMapper interface.
func NewImMapper(ctx context.Context) ImMapper {
	canCtx, terminate := context.WithCancel(ctx)
	retMap := &ImMap{canCtx, make(chan *ImPack, 1), make(chan *ImPack, 1), make(chan struct{}, 1), nil}

	retMap.stopMap = terminate

	go retMap.RunLoop()
	return retMap
}

func (imap *ImMap) Stop() {
	imap.stopMap()
	<-imap.done
}

// RunLoop is the ImMapper's map requests processing loop.
func (imap *ImMap) RunLoop() {

	pages := make(map[interface{}]interface{})
	for {
		select {
		case <-imap.Done():
			imap.done <- struct{}{}
			return
		case adder := <-imap.addChan:
			if _, exists := pages[adder.key]; exists {
				adder.ret <- fmt.Errorf("key exists")
				continue
			}
			pages[adder.key] = adder.value
			adder.ret <- nil
		case checker := <-imap.checkChan:
			if value, exists := pages[checker.key]; exists {
				checker.ret <- value
			} else {
				checker.ret <- nil
			}
		}
	}

}

// Add method allows one to add new keys.
// Returns error.
func (imap *ImMap) Add(key, value interface{}) error {
	iPack := &ImPack{key, value, make(chan interface{}, 1)}
	imap.addChan <- iPack

	val := <-iPack.ret
	if val == nil {
		return nil

	}
	if erval, ok := val.(error); ok {
		return errors.Wrap(erval, "Key Addition failed")
	}
	panic("panic in Add")
}

// Exists method allows to check and return the key.
func (imap *ImMap) Exists(key interface{}) (interface{}, bool) {
	iPack := &ImPack{key, nil, make(chan interface{}, 1)}
	imap.checkChan <- iPack
	val := <-iPack.ret

	if val == nil {
		return nil, false
	}
	return val, true
}
