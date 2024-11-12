VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
LDFLAGS = -X main.version=${VERSION}

build:
	go build -ldflags '$(LDFLAGS)' cmd/cosmos-proposals-checker.go

install:
	go install -ldflags '$(LDFLAGS)' cmd/cosmos-proposals-checker.go

lint:
	golangci-lint run --fix ./...

test:
	go test -coverprofile cover.out -coverpkg ./... -v ./...

coverage:
	go tool cover -html=cover.out
