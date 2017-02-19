package immap

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var testcmIntfs = []struct {
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

func TestAddExistsConImut(t *testing.T) {

	var ok bool

	mapper, cFunc := NewContextImutMapper(context.Background())
	defer cFunc()

	for _, kv := range testcmIntfs {
		mapper.Add(kv.key, kv.value)
	}

	for _, kv := range testcmIntfs {
		if _, ok, _ = mapper.Exists(kv.key); !ok {
			t.Fatal("Key addition failed, does not exist")
		}
	}

}

func TestAddExistsConImutConc(t *testing.T) {

	var ok bool
	var wg sync.WaitGroup

	mapper, cFunc := NewContextImutMapper(context.Background())
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
		for ind, _ := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			if _, ok, _ = mapper.Exists(ind); !ok {
				t.Fatal("Key addition failed, does not exist")
			}
		}
	}()

	wg.Wait()
}

func TestAddExistsConImutConcRand(t *testing.T) {

	var ok bool
	var wg sync.WaitGroup

	mapper, cFunc := NewContextImutMapper(context.Background())
	defer cFunc()

	count := 10000
	limit := 50

	randChan := make(chan int, count)
	randChan2 := make(chan int, count)
	go randGenChan(randChan, count, limit)
	go randGenChan(randChan2, count, limit)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for key := range randChan {
			mapper.Add(key, struct{}{})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for key := range randChan2 {
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
