package immap

import (
	"context"
	"testing"
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