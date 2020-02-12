.PHONY: test dep test-race vet lint analyse
export GO111MODULE = on


dep:
	go mod download

test: dep
	go test -v

test-race: dep
	go test -race -v

vet: lint
	go vet .

analyse: vet lint
