# Initialises the project after cloned
init:
	cp ./assets/certs/test.server.crt ./assets/certs/server.crt
	cp ./assets/certs/test.server.key ./assets/certs/server.key

certs:
	chmod +x ./scripts/generate.certs.sh
	mv ./server.crt ./assets/certs/server.crt
	mv ./server.key ./assets/certs/server.key

all: static-check test run

static-check:
	golangci-lint --tests=false run

## Cleanup package and detect then fix common linters warnings
lint:
	go mod tidy
	golangci-lint --fix --tests=false --timeout=2m30s run

## Run unit tests
test-unit:
	go test -v ./... -count=1

test-cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## Run lint and test-unit commands
test: lint test-unit

## Run test command and then build the executable
build: test
	go build -o bin/demo-rest-api-server -a -ldflags "-extldflags '-static' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" .

## Run lint and test-unit commands
run-local: init
	golangci-lint --tests=false run --fast
	go run -ldflags "-X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" main.go

## Build and run the app container
run-docker:
	docker-compose build && docker-compose up app

# Format codebase
format:
	gofumpt -l -w .