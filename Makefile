GO ?= go
COVERAGE_OUT = coverage.out

.PHONY: test test-norace

test:
	$(GO) test ./... -count=1 -race -coverprofile=$(COVERAGE_OUT)

test-norace:
	$(GO) test ./... -count=1

