.PHONY: test dep

# Credit: https://github.com/rightscale/go-boilerplate/blob/master/Makefile
DEPEND=golang.org/x/tools/cmd/cover github.com/Masterminds/glide github.com/golang/lint/golint

dep:
	go get -u $(DEPEND)
	glide install

test: dep
	go test -v

