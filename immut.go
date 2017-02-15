package immap

import (
	"context"
	"log"
)

const (
	DELETED = "marked as deleted"
)

// NewImutMapper returns a new instance of implementing ImutMapper interface.
func NewImutMapper(ctx context.Context) (ImutMapper, context.CancelFunc) {
	canCtx, terminate := context.WithCancel(ctx)
	iPack := &ImutMap{canCtx, make(chan *mapPack, 1)}

	go iPack.runLoop()
	return iPack, terminate
}

// RunLoop is the ImutMapper's map requests processing loop.
func (imap *ImutMap) runLoop() {

	pageList := make([]IntfMap, 0)

	for {
	SelAgain:
		select {
		case <-imap.Done():
			return
		case opMsg := <-imap.cChan:
			switch opMsg.op {
			case ADD_KEY:

				for counter := 0; counter <= len(pageList)-1; counter++ {
					pages := pageList[counter]
					value, exists := pages[opMsg.key]
					if !exists || value == DELETED {
						pages[opMsg.key] = opMsg.value
						opMsg.ret <- retPack{nil, pages}
						break SelAgain
					}
				}
				pageList = append(pageList, make(IntfMap))
				lpage := pageList[len(pageList)-1]
				lpage[opMsg.key] = opMsg.value
				opMsg.ret <- retPack{nil, lpage}
			case CHECK_KEY:
				for counter := len(pageList) - 1; counter >= 0; counter-- {
					pages := pageList[counter]
					if value, exists := pages[opMsg.key]; exists {
						if value == DELETED {
							opMsg.ret <- retPack{nil, nil}
							break SelAgain
						} else {
							opMsg.ret <- retPack{value, pages}
							break SelAgain
						}
					}
				}
				opMsg.ret <- retPack{nil, nil}

			case DEL_KEY:
				for counter := len(pageList) - 1; counter >= 0; counter-- {
					pages := pageList[counter]
					value, exists := pages[opMsg.key]
					if exists {
						if value == DELETED {
							// Do nothing
							log.Printf("Key %+v already deleted", opMsg.key)

						} else {
							pages[opMsg.key] = DELETED
							break SelAgain
						}
					}
				}

			}
		}
	}

}

// Add method allows one to add new keys.
// Returns error.
func (imap *ImutMap) Add(key, value interface{}) (IntfMap, error) {
	iPack := &mapPack{ADD_KEY, key, value, make(chan retPack, 1)}
	imap.cChan <- iPack

	pack := <-iPack.ret

	return pack.mapRef, nil
}

// Exists method allows to check and return the key.
func (imap *ImutMap) Exists(key interface{}) (interface{}, bool, IntfMap) {
	iPack := &mapPack{CHECK_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false, nil
	}
	return val.value, true, val.mapRef
}

func (imap *ImutMap) Delete(key interface{}) {
	iPack := &mapPack{DEL_KEY, key, nil, nil}
	imap.cChan <- iPack
}
