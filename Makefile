PROGRAM := screw
SOURCES := $(shell find ./cmd ./internal -type f -name '*.go')
CMD_DIR := cmd
BIN_DIR := .bin

$(BIN_DIR)/$(PROGRAM): $(SOURCES)
	go build -o ./$(BIN_DIR)/$(PROGRAM) ./$(CMD_DIR)/$(PROGRAM)

.PHONY: deps
deps:
	dep ensure

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install ./...

.PHONY: uninstall
uninstall:
	go clean -i ./$(CMD_DIR)/$(PROGRAM)

.PHONY: clean
clean:
	rm -rf ./$(BIN_DIR)/

# vim: noexpandtab
