package immap

import (
	"context"
)

// NewImutMapper returns a new instance of implementing ImutMapper interface.
func NewImutMapper(ctx context.Context) ImutMapper {
	canCtx, terminate := context.WithCancel(ctx)
	retMap := &ImutMap{canCtx, make(chan *ImPack, 1), make(chan *ImPack, 1), make(chan struct{}, 1), nil}

	retMap.stopMap = terminate

	go retMap.RunImLoop()
	return retMap
}

type IntfMap map[interface{}]interface{}

func (imap *ImutMap) Stop() {
	imap.stopMap()
	<-imap.done
}

// RunLoop is the ImutMapper's map requests processing loop.
func (imap *ImutMap) RunImLoop() {

	pageList := make([]IntfMap, 0)
	var added bool

	for {
		added = false
		select {
		case <-imap.Done():
			imap.done <- struct{}{}
			return
		case adder := <-imap.addChan:

			for counter := 0; counter <= len(pageList)-1; counter++ {
				pages := pageList[counter]
				if _, exists := pages[adder.key]; !exists {
					pages[adder.key] = adder.value
					adder.ret <- retMap{nil, pages}
					added = true
					break
				}
			}
			if added == false {
				pageList = append(pageList, make(map[interface{}]interface{}))
				lpage := pageList[len(pageList)-1]
				lpage[adder.key] = adder.value
				adder.ret <- retMap{nil, lpage}
			}
		case checker := <-imap.checkChan:
			for counter := len(pageList) - 1; counter >= 0; counter-- {
				pages := pageList[counter]
				if value, exists := pages[checker.key]; exists {
					checker.ret <- retMap{value, pages}
					added = true
					break
				}
			}
			if added == false {
				checker.ret <- retMap{nil, nil}
			}

		}
	}

}

// Add method allows one to add new keys.
// Returns error.
func (imap *ImutMap) Add(key, value interface{}) (IntfMap, error) {
	iPack := &ImPack{key, value, make(chan retMap, 1)}
	imap.addChan <- iPack

	pack := <-iPack.ret

	return pack.mapRef, nil
}

// Exists method allows to check and return the key.
func (imap *ImutMap) Exists(key interface{}) (interface{}, bool, IntfMap) {
	iPack := &ImPack{key, nil, make(chan retMap, 1)}
	imap.checkChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false, nil
	}
	return val.value, true, val.mapRef
}
