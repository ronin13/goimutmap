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
	{make(chan struct{}), struct{}{}},
	{struct{ x int }{33}, "zz"},
	{[...]int{1, 3, 4}, []byte("aaaaa")},
}

func TestAddExists(t *testing.T) {

	var err error
	var ok bool

	mapper := NewcontextMapper(context.Background())
	defer mapper.Stop()

	for _, kv := range testIntfs {
		err = mapper.Add(kv.key, kv.value)
		if err != nil {
			t.Fatal("Addition failed")
		}
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

func TestAddExistsImut(t *testing.T) {

	var err error
	var ok bool

	mapper := NewImutMapper(context.Background())
	defer mapper.Stop()

	for _, kv := range testIntfs {
		_, err = mapper.Add(kv.key, kv.value)
		if err != nil {
			t.Fatal("Addition failed")
		}
	}

	for _, kv := range testIntfs {
		if _, ok, _ = mapper.Exists(kv.key); !ok {
			t.Fatal("Key addition failed, does not exist")
		}
	}

	for _, kv := range testIntfs {
		mapper.Delete(kv.key)
	}
}

func TestAddExistsConc(t *testing.T) {

	var err error
	var ok bool
	var wg sync.WaitGroup

	mapper := NewcontextMapper(context.Background())
	defer mapper.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ind, kv := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			err = mapper.Add(ind, kv)
			if err != nil {
				t.Fatal("Addition failed")
			}
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

func TestAddExistsImutConc(t *testing.T) {

	var err error
	var ok bool
	var wg sync.WaitGroup

	mapper := NewImutMapper(context.Background())
	defer mapper.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ind, kv := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			_, err = mapper.Add(ind, kv)
			if err != nil {
				t.Fatal("Addition failed")
			}
		}
	}()

	wg.Add(1)
	time.Sleep(2 * time.Second)
	go func() {
		defer wg.Done()
		for ind, _ := range [100]struct{}{} {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			if _, ok, _ = mapper.Exists(ind); !ok {
				t.Fatal("Key addition failed, does not exist")
			}
		}
	}()

	wg.Wait()
}
