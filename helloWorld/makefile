MAIN_PACKAGES ?= $(shell go list -f '{{ .ImportPath }}:{{ .Name }}' ./cmd/... | grep ":main" | grep -v -e "/examples/" | cut -d: -f1)
LDFLAGS ?= "-s -w" # for light version
GOBUILD ?= go build -ldflags=$(LDFLAGS)

clean:
	rm -rf bin

run:
	go run main.go

build: clean
	@echo "Building..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	$(GOBUILD) -o bin/ $(MAIN_PACKAGES)

test:
	go clean -testcache
	@echo "=== Running tests ==="
	go test -timeout 120s -failfast -race -coverprofile=test.cover ./...

	@echo "=== Show coverage ==="
	go tool cover -func=test.cover
	
	@echo "=== Cleaning up test.cover ==="
	rm -f test.cover