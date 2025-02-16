lint:
    golangci-lint run --fix

test: lint
    go test -v ./...
