.PHONY: clean deps docker build

EXECUTABLE ?= maze-publish
IMAGE ?= mikkeloscar/$(EXECUTABLE)

all: build

clean:
	go clean -i ./..

deps:
	go get -t

docker: build
	docker build --rm -t $(IMAGE) .

$(EXECUTABLE): $(wildcard *.go)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s"

build: $(EXECUTABLE)
