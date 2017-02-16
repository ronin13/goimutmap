package immap

import (
	"context"
	"log"
)

const (
	// DELETED is a delete marker used by ImutMap.
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

	mapList := make([]IntfMap, 0)

	for {
	SelAgain:
		select {
		case <-imap.Done():
			return
		case opMsg := <-imap.cChan:
			switch opMsg.op {
			case ADD_KEY:
				for counter := 0; counter <= len(mapList)-1; counter++ {
					indMap := mapList[counter]
					value, exists := indMap[opMsg.key]
					if !exists || (exists && value == DELETED) {
						indMap[opMsg.key] = opMsg.value
						opMsg.ret <- retPack{nil, indMap}
						break SelAgain
					}
				}
				mapList = append(mapList, make(IntfMap))
				lpage := mapList[len(mapList)-1]
				lpage[opMsg.key] = opMsg.value
				opMsg.ret <- retPack{nil, lpage}
			case CHECK_KEY:
				for counter := len(mapList) - 1; counter >= 0; counter-- {
					indMap := mapList[counter]
					if value, exists := indMap[opMsg.key]; exists {
						if value == DELETED {
							opMsg.ret <- retPack{nil, nil}
							break SelAgain
						} else {
							opMsg.ret <- retPack{value, indMap}
							break SelAgain
						}
					}
				}
				opMsg.ret <- retPack{nil, nil}
			case DEL_KEY:
				for counter := len(mapList) - 1; counter >= 0; counter-- {
					indMap := mapList[counter]
					value, exists := indMap[opMsg.key]
					if exists {
						if value == DELETED {
							// Do nothing
							log.Printf("Key %+v already deleted", opMsg.key)
							break SelAgain

						} else {
							indMap[opMsg.key] = DELETED
							break SelAgain
						}
					}
				}

			}
		}
	}

}

// Add method allows one to add new keys.
// Returns a reference to IntfMap
func (imap *ImutMap) Add(key, value interface{}) IntfMap {
	iPack := &mapPack{ADD_KEY, key, value, make(chan retPack, 1)}
	imap.cChan <- iPack

	pack := <-iPack.ret

	return pack.mapRef
}

// Exists method allows to check and return the key.
// Also, returns a reference to IntfMap.
func (imap *ImutMap) Exists(key interface{}) (interface{}, bool, IntfMap) {
	iPack := &mapPack{CHECK_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false, nil
	}
	return val.value, true, val.mapRef
}

// Delete method allows to delete keys.
func (imap *ImutMap) Delete(key interface{}) {
	iPack := &mapPack{DEL_KEY, key, nil, nil}
	imap.cChan <- iPack
}
