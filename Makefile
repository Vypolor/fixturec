MD_DIR ?= ./cmd/fixturec
BIN_DIR ?= ./bin
BIN_NAME ?= fixturec
BIN := $(BIN_DIR)/$(BIN_NAME)

GO ?= go
GOCMD := $(GO)

# Build the CLI binary
build: $(BIN_DIR) $(BIN)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN):
	@echo "Building $(BIN)..."
	$(GOCMD) build -o $(BIN) $(CMD_DIR)
