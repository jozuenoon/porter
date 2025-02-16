lint:
    golangci-lint run --fix

test:
    go test -v ./...

bin:
   CGO_ENABLED=0 go build -o bin/porter -ldflags="-s -w" cmd/main.go

clean:
    rm -rf bin

build: clean lint test bin
    docker build -t q88/porter:latest -f deployment/Dockerfile .

upload:
    # NOTE: Please replace with large ports json file to test memory usage.
    curl -T testdata/ports.json -X POST -H "Content-Type: application/json" http://localhost:8080/api/v1/ports

get:
    curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/port?unloc=AEAJM

run:
    docker run -it -p 8080:8080 --memory=200m q88/porter:latest
