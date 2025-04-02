build:
    @mkdir -p bin
    @go build -o ./bin/cdkpw ./cmd/...

test:
    @go test ./cmd/... -coverprofile=coverage.out

install:
    @go install ./cmd/...
