package immap

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var testmIntfs = []struct {
	key   interface{}
	value interface{}
}{
	{"abc", "123"},
	{make(chan struct{}), struct{}{}},
	{struct{ x int }{33}, "zz"},
	{[...]int{1, 3, 4}, []byte("aaaaa")},
	{[...]int{1, 3, 4}, []byte("aaaaa")},
	{[...]int{1, 3, 4}, []byte("aaaaa")},
}

func TestAddExistsImut(t *testing.T) {

	var ok bool

	mapper, cFunc := NewImutMapper(context.Background())
	defer cFunc()

	for _, kv := range testmIntfs {
		mapper.Add(kv.key, kv.value)
	}

	for _, kv := range testmIntfs {
		if _, ok, _ = mapper.Exists(kv.key); !ok {
			t.Fatal("Key addition failed, does not exist")
		}
	}

	for _, kv := range testmIntfs {
		mapper.Delete(kv.key)
	}
}

func TestAddExistsImutConc(t *testing.T) {

	var ok bool
	var wg sync.WaitGroup

	mapper, cFunc := NewImutMapper(context.Background())
	defer cFunc()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ind, kv := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(1)) * time.Millisecond)
			mapper.Add(ind, kv)
		}
	}()

	wg.Add(1)
	time.Sleep(2 * time.Second)
	go func() {
		defer wg.Done()
		for ind := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			if _, ok, _ = mapper.Exists(ind); !ok {
				t.Fatal("Key addition failed, does not exist")
			}
		}
	}()

	wg.Wait()
}

func randGenChan(rchan chan int, count, limit int) {
	for x := 0; x < count; x++ {
		rchan <- rand.Intn(limit)
	}
	close(rchan)
}

func TestAddExistsImutConcRand(t *testing.T) {

	var ok bool
	var wg sync.WaitGroup

	mapper, cFunc := NewImutMapper(context.Background())
	defer cFunc()

	count := 10000
	limit := 3

	randChan := make(chan int, count)
	go randGenChan(randChan, count, limit)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for key := range randChan {
			mapper.Add(key, struct{}{})
		}
	}()

	wg.Add(1)
	time.Sleep(2 * time.Second)
	go func() {
		defer wg.Done()
		for key := range randChan {
			if _, ok, _ = mapper.Exists(key); !ok {
				t.Fatal("Key addition failed, does not exist")
			}
		}
	}()

	wg.Wait()
}
