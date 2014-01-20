.PHONY: all clean deps fmt check test build

all: fmt deps check test build

clean:
	go clean

deps:
	go get -d -v
	go get -u github.com/golang/lint/golint # frequently updated, so -u

fmt:
	gofmt -w .

# lint and vet both return success (0) on error. make them error and report
check: deps
	go tool vet . 2>&1 | wc -l | { grep 0 || { go tool vet . && false; }; }
	if find . -name '*.go' | xargs golint | grep ":"; then false; else true; fi

test: deps
	go test -v

build: deps
	go build -v
