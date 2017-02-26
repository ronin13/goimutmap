package immap

import (
	"context"
	"reflect"
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

	bookMap := make(map[interface{}][]int)
	var lastInd, lastMapInd int
	var iPoint int

	for {
		select {
		case <-imap.Done():
			return
		case opMsg := <-imap.cChan:
			startSlice, exists := bookMap[opMsg.key]

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

				bookMap[opMsg.key] = append(bookMap[opMsg.key], lastInd)

				opMsg.ret <- retPack{nil, lpage}
			case CHECK_KEY:
				intMap := mapList[lastInd]
				value, _ := intMap[opMsg.key]
				if value == DELETED {
					panic("DELETED value should not be here")
				}
				opMsg.ret <- retPack{value, intMap}
			case DEL_KEY:
				intMap := mapList[lastInd]
				value, exists := intMap[opMsg.key]
				if exists {
					if value == DELETED {
						panic("DELETED value should not be here")

					} else {
						intMap[opMsg.key] = DELETED
						bookMap[opMsg.key] = bookMap[opMsg.key][:len(bookMap[opMsg.key])-1]
					}
				}

			}
		}
	}

}

// Add method allows one to add new keys.
// Returns a reference to IntfMap
func (imap *ImutMap) Add(key, value interface{}) IntfMap {

	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	iPack := &mapPack{ADD_KEY, key, value, make(chan retPack, 1)}
	imap.cChan <- iPack

	pack := <-iPack.ret

	return (pack.mapRef).(IntfMap)
}

// Exists method allows to check and return the key.
// Also, returns a reference to IntfMap.
func (imap *ImutMap) Exists(key interface{}) (interface{}, bool, IntfMap) {
	iPack := &mapPack{CHECK_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	val := <-iPack.ret

	tmap := (val.mapRef).(IntfMap)

	_, exists := tmap[key]

	return val.value, exists, tmap
}

// Delete method allows to delete keys.
func (imap *ImutMap) Delete(key interface{}) {
	iPack := &mapPack{DEL_KEY, key, nil, nil}
	imap.cChan <- iPack
}
