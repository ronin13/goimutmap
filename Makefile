.PHONY: test dep test-race vet lint analyse

# Credit: https://github.com/rightscale/go-boilerplate/blob/master/Makefile
DEPEND=golang.org/x/tools/cmd/cover  golang.org/x/lint@latest

dep:
	go mod download
	go get -u $(DEPEND)

test: dep
	go test -v

test-race: dep
	go test -race -v

vet: lint
	go vet .

lint:
	golint .

analyse: vet lint
