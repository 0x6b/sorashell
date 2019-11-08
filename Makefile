GO=go
BIN=soracom-shell
BIN_DOWNLOAD=download-assets
ASSETS = assets/en.yaml assets/ja.yaml assets/soracom-api.en.yaml assets/soracom-api.ja.yaml

$(BIN): clean $(ASSETS)
	$(GO) build -o $(BIN) ./cmd/shell

$(ASSETS): $(BIN_DOWNLOAD)
	./$(BIN_DOWNLOAD)

$(BIN_DOWNLOAD):
	$(GO) build -o $(BIN_DOWNLOAD) ./scripts/download

test: $(BIN)
	$(GO) test -v

clean:
	rm -f $(BIN) $(BIN_DOWNLOAD) $(ASSETS)
	$(GO) clean

.PHONY: test clean
