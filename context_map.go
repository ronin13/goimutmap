package immap

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

// NewcontextMapper returns a new instance of implementing contextMapper interface.
func NewcontextMapper(ctx context.Context) (contextMapper, context.CancelFunc) {
	canCtx, terminate := context.WithCancel(ctx)
	cPack := &contextMap{canCtx, make(chan *mapPack, 1), make(chan *mapPack, 1), make(chan *mapPack, 1)}

	go cPack.runLoop()
	return cPack, terminate
}

// RunLoop is the contextMapper's map requests processing loop.
func (imap *contextMap) runLoop() {

	pages := make(IntfMap)
	for {
		select {
		case <-imap.Done():
			return
		case adder := <-imap.addChan:
			if _, exists := pages[adder.key]; exists {
				adder.ret <- retPack{fmt.Errorf("key exists"), nil}
				continue
			}
			pages[adder.key] = adder.value
			adder.ret <- retPack{nil, nil}
		case checker := <-imap.checkChan:
			if value, exists := pages[checker.key]; exists {
				checker.ret <- retPack{value, nil}
			} else {
				checker.ret <- retPack{nil, nil}
			}
		case deler := <-imap.delChan:
			delete(pages, deler.key)
			deler.ret <- retPack{nil, nil}

		}
	}

}

// Add method allows one to add new keys.
// Returns error.
func (imap *contextMap) Add(key, value interface{}) error {
	iPack := &mapPack{key, value, make(chan retPack, 1)}
	imap.addChan <- iPack

	val := <-iPack.ret
	if val.value == nil {
		return nil

	}
	if erval, ok := val.value.(error); ok {
		return errors.Wrap(erval, "Key Addition failed")
	}
	panic("panic in Add")
}

// Exists method allows to check and return the key.
func (imap *contextMap) Exists(key interface{}) (interface{}, bool) {
	iPack := &mapPack{key, nil, make(chan retPack, 1)}
	imap.checkChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false
	}
	return val.value, true
}

func (imap *contextMap) Delete(key interface{}) {
	iPack := &mapPack{key, nil, make(chan retPack, 1)}
	imap.checkChan <- iPack
	_ = <-iPack.ret
}
