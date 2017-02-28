[![Sourcegraph](https://sourcegraph.com/github.com/ronin13/goimutmap/-/badge.svg)](https://sourcegraph.com/github.com/ronin13/goimutmap?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/ronin13/goimutmap)](https://goreportcard.com/report/github.com/ronin13/goimutmap)
[![Build Status](https://travis-ci.org/ronin13/goimutmap.svg?branch=master)](https://travis-ci.org/ronin13/goimutmap)

[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/ronin13/goimutmap/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/ronin13/goimutmap?status.svg)](https://godoc.org/github.com/ronin13/goimutmap)


# goimutmap

## Introduction

This library provides:

1) A lockless Golang map implementing the Context interface.  (ContextMapper)

2) An immutable multi-versioned map built on top of the lockless map, similarly implementing the Context interface.  (ImmutMapper)

### Supports: (with similar semantics as that of golang map)

a) Add

b) Exists

c) Delete

d) Iterate

### ContextMapper

```
type IterMap struct {
	key, value interface{}
}

type ContextMapper interface {
	Exists(interface{}) (interface{}, bool)
	Add(interface{}, interface{}) interface{}
	Delete(interface{})
	Iterate() <-chan IterMap
}
```

As can be seen from above, it implements an interface similar to that of regular map.

### ImmutMapper

```
type IntfMap map[interface{}]interface{}

type ImutMapper interface {
	Exists(interface{}) (interface{}, bool, IntfMap)
	Add(interface{}, interface{}) IntfMap
	Delete(interface{})
}
```

ImmutMapper implements similar interface, except it returns a 3rd value 
which is a `snapshot` of the map into which the operation was done.

Both ContextMapper and ImmutMapper  encapsulate context.Context:

```
type baseMap struct {
	context.Context
	// Other internal fields
}

```

and provide constructors such as:

```
NewcontextMapper(ctx context.Context) (ContextMapper, context.CancelFunc)
```

and 

```
NewImutMapper(ctx context.Context) (ImutMapper, context.CancelFunc)
```

and finally,

```
type ContextImutMapper interface {
	Exists(interface{}) (interface{}, bool, ContextMapper)
	Add(interface{}, interface{}) ContextMapper
}
```

which combines both ContextMapper and ImmutMapper.



#### Note
The key and values inserted can by of any type and heterogenous.

Please refer to [godoc](https://godoc.org/github.com/ronin13/goimutmap) for more details.

## Used by
* http://github.com/ronin13/dotler : Multiple crawler goroutines use ContextMapper to avoid duplicate crawling and for in-memory graph. 

## Examples:
* Usage: https://github.com/ronin13/dotler/blob/master/wire/nodemap.go#L11
