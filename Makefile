build:
	go build cmd/cosmos-proposals-checker.go

install:
	go install cmd/cosmos-proposals-checker.go

lint:
	golangci-lint run --fix ./...

test:
	go test -coverprofile cover.out -v ./...
