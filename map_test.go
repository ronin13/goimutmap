package immap

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var testIntfs = []struct {
	key   interface{}
	value interface{}
}{
	{"abc", "123"},
	{"abc", "123"},
	{make(chan struct{}), struct{}{}},
	{struct{ x int }{33}, "zz"},
	{[...]int{1, 3, 4}, []byte("aaaaa")},
}

func TestAddExists(t *testing.T) {

	var ok bool

	mapper, cFunc := NewcontextMapper(context.Background())
	defer cFunc()

	for _, kv := range testIntfs {
		mapper.Add(kv.key, kv.value)
	}

	for _, kv := range testIntfs {
		if _, ok = mapper.Exists(kv.key); !ok {
			t.Fatal("Key addition failed, does not exist")
		}
	}

	for _, kv := range testIntfs {
		mapper.Delete(kv.key)
	}
}

func TestAddExistsConc(t *testing.T) {

	var ok bool
	var wg sync.WaitGroup

	mapper, cFunc := NewcontextMapper(context.Background())
	defer cFunc()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ind, kv := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			mapper.Add(ind, kv)
		}
	}()

	wg.Add(1)
	time.Sleep(2 * time.Second)
	go func() {
		defer wg.Done()
		for ind, _ := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			if _, ok = mapper.Exists(ind); !ok {
				t.Fatal("Key addition failed, does not exist")
			}
		}
	}()

	wg.Wait()
}
