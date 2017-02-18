package immap

import (
	"context"
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

	intMap := make(map[interface{}][]int)
	var lastInd, lastMapInd int
	var iPoint int

	for {
	SelAgain:
		select {
		case <-imap.Done():
			return
		case opMsg := <-imap.cChan:
			startSlice, exists := intMap[opMsg.key]

			if !exists {
				lastInd = -1
			} else {
				lastInd = startSlice[len(startSlice)-1]
			}
			switch opMsg.op {
			case ADD_KEY:

				expand := false

				if len(mapList) > 0 {
					lastMapInd = len(mapList) - 1
				} else {
					expand = true
				}

				lastInd = lastInd + 1

				if lastInd > lastMapInd || expand {
					mapList = append(mapList, make(IntfMap))
					iPoint = len(mapList) - 1
				} else {
					iPoint = lastInd
				}

				lpage := mapList[iPoint]
				lpage[opMsg.key] = opMsg.value

				intMap[opMsg.key] = append(intMap[opMsg.key], lastInd)

				opMsg.ret <- retPack{nil, lpage}
			case CHECK_KEY:
				indMap := mapList[lastInd]
				if value, exists := indMap[opMsg.key]; exists {
					if value == DELETED {
						panic("DELETED value should not be here")
					} else {
						opMsg.ret <- retPack{value, indMap}
						break SelAgain
					}
				}
				opMsg.ret <- retPack{nil, nil}
			case DEL_KEY:
				indMap := mapList[lastInd]
				value, exists := indMap[opMsg.key]
				if exists {
					if value == DELETED {
						panic("DELETED value should not be here")

					} else {
						indMap[opMsg.key] = DELETED
						intMap[opMsg.key] = intMap[opMsg.key][:len(intMap[opMsg.key])-1]
						break SelAgain
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
