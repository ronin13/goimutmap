package immap

import (
	"context"
	"reflect"
)

const (
	// DELETED is a delete marker used by ConImutMap.
	CDELETED = "marked as deleted"
)

// NewConImutMapper returns a new instance of implementing ConImutMapper interface.
func NewContextImutMapper(ctx context.Context) (ContextImutMapper, context.CancelFunc) {
	canCtx, terminate := context.WithCancel(ctx)
	iPack := &ConImutMap{canCtx, make(chan *mapPack, 1)}

	go iPack.runLoop()
	return iPack, terminate
}

func (imap *ConImutMap) runLoop() {

	mapList := make([]ContextMapper, 0)

	bookMap, _ := NewcontextMapper(imap.Context)
	var lastInd, lastMapInd int
	var iPoint int
	var tMap ContextMapper
	var ok bool

	for {
		select {
		case <-imap.Done():
			return

		case opMsg := <-imap.cChan:
			startSlice, exists := bookMap.Exists(opMsg.key)
			iList := []int{}

			if !exists {
				lastInd = -1
			} else {
				iList, ok = startSlice.([]int)

				if !ok {
					panic("Value of bookMap needs to be a []int")
				}
				lastInd = iList[len(iList)-1]
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

				if expand || lastInd > lastMapInd {
					tMap, _ = NewcontextMapper(imap.Context)
					mapList = append(mapList, tMap)
					iPoint = len(mapList) - 1
				} else {
					iPoint = lastInd
				}

				lpage := mapList[iPoint]
				lpage.Add(opMsg.key, opMsg.value)

				_, exists := bookMap.Exists(opMsg.key)

				if exists {
					bookMap.Add(opMsg.key, append(iList, lastInd))
				} else {
					bookMap.Add(opMsg.key, append([]int(nil), lastInd))
				}

				opMsg.ret <- retPack{nil, lpage}

			case CHECK_KEY:

				intMap := mapList[lastInd]
				if value, exists := intMap.Exists(opMsg.key); exists {
					if value == CDELETED {
						panic("CDELETED value should not be here")
					} else {
						opMsg.ret <- retPack{value, intMap}
					}
				} else {
					opMsg.ret <- retPack{nil, nil}
				}

			}

		}
	}
}

// Add method allows one to add new keys.
// Returns a reference to ContextMapper
func (imap *ConImutMap) Add(key, value interface{}) ContextMapper {

	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	iPack := &mapPack{ADD_KEY, key, value, make(chan retPack, 1)}
	imap.cChan <- iPack

	pack := <-iPack.ret

	return (pack.mapRef).(ContextMapper)

}

// Exists method allows to check and return the key.
// Also, returns a reference to ContextMapper.
func (imap *ConImutMap) Exists(key interface{}) (interface{}, bool, ContextMapper) {

	iPack := &mapPack{CHECK_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false, nil
	}
	return val.value, true, (val.mapRef).(ContextMapper)

}
