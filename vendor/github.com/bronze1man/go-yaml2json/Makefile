.PHONY: all deps build test cov install doc clean

all: test install

deps:
	go get -d -v -t ./...

build: deps
	go build ./...

test: deps
	go test -test.v ./...

cov: deps
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

install: deps
	go install ./...

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

clean:
	go clean -i ./...
