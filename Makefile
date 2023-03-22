GO_TEST_FILES := $(shell find . -name '*_test.go' | grep -v /vendor/)

.PHONY: test
test:
	@go test -v -cover -race ./...
	
.PHONY: test-coverprofile
test-coverprofile:
	@go test -v -cover -coverprofile=coverage.out -covermode atomic -race ./...
	
.PHONY: deps
deps:
	@go mod download

.PHONY: build
build:
	mkdir -p bin
	@CGO_ENABLED=0 go build -o bin/server ./cmd/server/...

.PHONY: fmt
fmt:
	@goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: fmt-check
fmt-check: fmt
	@DIFF=$$(git diff) && \
	if [ -n "$$DIFF" ]; then \
		echo "$$DIFF" && echo "run 'make fmt' and commit the changes" && exit 1; \
	fi

.PHONY: lint
lint:
	@golangci-lint run --timeout '10m0s'

coverage.out: ${GO_TEST_FILES}
	@go test -coverprofile=coverage.out ./...

.PHONY: show-coverage
show-coverage: coverage.out
	@go tool cover -html=coverage.out

