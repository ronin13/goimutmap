package immap

import (
	"context"
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

type IntfMap map[interface{}]interface{}

// RunLoop is the ImMapper's map requests processing loop.
func (imap *ImMap) RunLoop() {

	pageList := make([]IntfMap, 0)
	var counter int
	var added bool

	for {
		select {
		case <-imap.Done():
			imap.done <- struct{}{}
			return
		case adder := <-imap.addChan:

			added = false
			for counter = 0; counter <= len(pageList)-1; counter++ {
				pages := pageList[counter]
				if _, exists := pages[adder.key]; !exists {
					pages[adder.key] = adder.value
					adder.ret <- nil
					added = true
					break
				}
			}
			if added == false {
				pageList = append(pageList, make(map[interface{}]interface{}))
				pageList[len(pageList)-1][adder.key] = adder.value
				adder.ret <- nil
			}
		case checker := <-imap.checkChan:
			counter = 0
			for counter = len(pageList) - 1; counter >= 0; counter-- {
				pages := pageList[counter]
				if value, exists := pages[checker.key]; exists {
					checker.ret <- value
					break
				}
			}
			if counter < 0 {
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
