
[![Go Report Card](https://goreportcard.com/badge/github.com/ronin13/goimutmap)](https://goreportcard.com/report/github.com/ronin13/goimutmap)
[![Build Status](https://travis-ci.org/ronin13/goimutmap.svg?branch=master)](https://travis-ci.org/ronin13/goimutmap)

[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/ronin13/goimutmap/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/ronin13/goimutmap?status.svg)](https://godoc.org/github.com/ronin13/goimutmap)


# goimutmap

## Introduction

This library provides:

1) A lockless Golang map implementing the Context interface. 

2) An immutable multi-versioned map built on top of the lockless map, similarly implementing the Context interface. (This is a Work In Progress)

Supports: (with similar semantics as that of golang map)

a) Add

b) Exists

c) Delete

Please refer to [godoc](https://godoc.org/github.com/ronin13/goimutmap) for details.

## Used by
* http://github.com/ronin13/dotler : Multiple crawler goroutines use this map to avoid duplicate crawling and for in-memory graph. 

## Examples:
* Usage: https://github.com/ronin13/dotler/blob/master/wire/nodemap.go#L11
