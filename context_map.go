package immap

import (
	"context"
)

// NewcontextMapper returns a new instance of implementing contextMapper interface.
func NewcontextMapper(ctx context.Context) (contextMapper, context.CancelFunc) {
	canCtx, terminate := context.WithCancel(ctx)
	cPack := &contextMap{canCtx, make(chan *mapPack, 1)}

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
		case opMsg := <-imap.cChan:
			switch opMsg.op {
			case ADD_KEY:
				if value, exists := pages[opMsg.key]; exists {
					opMsg.ret <- retPack{value, nil}
				} else {
					pages[opMsg.key] = opMsg.value
					opMsg.ret <- retPack{nil, nil}
				}
			case CHECK_KEY:
				if value, exists := pages[opMsg.key]; exists {
					opMsg.ret <- retPack{value, nil}
				} else {
					opMsg.ret <- retPack{nil, nil}
				}
			case DEL_KEY:
				delete(pages, opMsg.key)
				opMsg.ret <- retPack{nil, nil}
			}

		}
	}

}

// Add method allows one to add new keys.
func (imap *contextMap) Add(key, value interface{}) interface{} {
	iPack := &mapPack{ADD_KEY, key, value, make(chan retPack, 1)}
	imap.cChan <- iPack

	val := <-iPack.ret
	return val.value
}

// Exists method allows to check and return the key.
func (imap *contextMap) Exists(key interface{}) (interface{}, bool) {
	iPack := &mapPack{CHECK_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	val := <-iPack.ret

	if val.value == nil {
		return nil, false
	}
	return val.value, true
}

func (imap *contextMap) Delete(key interface{}) {
	iPack := &mapPack{DEL_KEY, key, nil, make(chan retPack, 1)}
	imap.cChan <- iPack
	_ = <-iPack.ret
}
