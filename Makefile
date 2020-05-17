### Required tools
GOTOOLS_CHECK = golangci-lint

all: fmt ensure-deps linter test

### Testing
test:
	go test ./... -covermode=atomic -coverpkg=./... -count=1 -race

test-cover:
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	go tool cover -html=/tmp/coverage.out

test-integration:
	go test -tags integration ./... -covermode=atomic -coverpkg=./... -count=1 -race

### Formatting, linting, and deps
fmt:
	go fmt ./...

ensure-deps:
	@echo "==> Running go mod tidy"
	go mod download
	go mod tidy

linter:
	@echo "==> Running linter"
	golangci-lint run ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: check_tools test test-cover fmt linter ensure-deps