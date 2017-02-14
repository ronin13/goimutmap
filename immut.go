package immap

import (
	"context"
)

const (
	DELETED = "marked as deleted"
)

// NewImutMapper returns a new instance of implementing ImutMapper interface.
func NewImutMapper(ctx context.Context) (ImutMapper, context.CancelFunc) {
	canCtx, terminate := context.WithCancel(ctx)
	retPack := &ImutMap{canCtx, make(chan *mapPack, 1), make(chan *mapPack, 1), make(chan *mapPack, 1)}

	go retPack.runLoop()
	return retPack, terminate
}

// RunLoop is the ImutMapper's map requests processing loop.
func (imap *ImutMap) runLoop() {

	pageList := make([]IntfMap, 0)
	var added bool

	for {
		added = false
		select {
		case <-imap.Done():
			return
		case adder := <-imap.addChan:

			for counter := 0; counter <= len(pageList)-1; counter++ {
				pages := pageList[counter]
				if _, exists := pages[adder.key]; !exists {
					pages[adder.key] = adder.value
					adder.ret <- retPack{nil, pages}
					added = true
					break
				}
			}
			if added == false {
				pageList = append(pageList, make(map[interface{}]interface{}))
				lpage := pageList[len(pageList)-1]
				lpage[adder.key] = adder.value
				adder.ret <- retPack{nil, lpage}
			}
		case checker := <-imap.checkChan:
			for counter := len(pageList) - 1; counter >= 0; counter-- {
				pages := pageList[counter]
				if value, exists := pages[checker.key]; exists {
					if value == DELETED {
						checker.ret <- retPack{nil, nil}
						added = true
						break
					} else {
						checker.ret <- retPack{value, pages}
						added = true
						break
					}
				}
			}
			if added == false {
				checker.ret <- retPack{nil, nil}
			}

		case deler := <-imap.delChan:
			for counter := len(pageList) - 1; counter >= 0; counter-- {
				pages := pageList[counter]
				if _, exists := pages[deler.key]; exists {
					pages[deler.key] = DELETED
					deler.ret <- retPack{nil, nil}
					break
				}
			}

		}
	}

}

// Add method allows one to add new keys.
// Returns error.
func (imap *ImutMap) Add(key, value interface{}) (IntfMap, error) {
	iPack := &mapPack{key, value, make(chan retPack, 1)}
	imap.addChan <- iPack

	pack := <-iPack.ret

	return pack.mapRef, nil
}

// Exists method allows to check and return the key.
func (imap *ImutMap) Exists(key interface{}) (interface{}, bool, IntfMap) {
	iPack := &mapPack{key, nil, make(chan retPack, 1)}
	imap.checkChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false, nil
	}
	return val.value, true, val.mapRef
}

func (imap *ImutMap) Delete(key interface{}) {
	iPack := &mapPack{key, nil, make(chan retPack, 1)}
	imap.checkChan <- iPack
	_ = <-iPack.ret
}
