PATHMAIN ?= ./cmd/server/main.go

dev:
    go run $(PATHMAIN)

.PHONY: dev
