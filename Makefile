GO ?= go

.PHONY: test test-norace

test:
	$(GO) test ./... -count=1 -race

test-norace:
	$(GO) test ./... -count=1


