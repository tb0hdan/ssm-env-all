GOPATH=$(shell pwd)

all:
	@go get -v -d .
	@go build .
