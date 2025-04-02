build:
    @mkdir -p bin
    @go build -o ./bin/cdkpw ./cmd/cdkpw/...

test:
    @go test ./cmd/cdkpw/... -coverprofile=coverage.out

install:
    @go install ./cmd/cdkpw/...
