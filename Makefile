REPO = github.com/chain710/copilot
TOOLS_DIR := $(shell pwd)/tools
export PATH := $(PATH):$(TOOLS_DIR)
PIGEON=$(TOOLS_DIR)/pigeon

all: mod check bin
bin: mod
	go build -o bin/copilot $(REPO)
test: mod check
check: fmt vet
fmt:
	go fmt ./...
vet:
	go vet ./...
mod:
	go mod download
generate: $(PIGEON)
	go generate ./...
$(PIGEON):
	GOBIN=$(TOOLS_DIR) go install github.com/mna/pigeon@v1.1.0