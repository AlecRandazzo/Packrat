# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=gofor-collector.exe

default: build
all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/gofor-collector/main.go
test:
		$(GOTEST) -race -v .