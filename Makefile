GO=go
BIN=soracom-shell
SRC= *.go cmd/shell/*.go
BIN_DOWNLOAD=download-assets
ASSETS = assets/en.yaml assets/ja.yaml assets/soracom-api.en.yaml assets/soracom-api.ja.yaml
STATIK = statik/statik.go

$(BIN): $(SRC) $(STATIK)
	$(GO) build -o $(BIN) ./cmd/shell

$(STATIK): $(ASSETS)
	statik -src=assets

$(ASSETS): $(BIN_DOWNLOAD)
	./$(BIN_DOWNLOAD)

$(BIN_DOWNLOAD):
	$(GO) build -o $(BIN_DOWNLOAD) ./scripts/download

test: $(BIN)
	$(GO) test -v

clean:
	rm -fr $(BIN) $(BIN_DOWNLOAD) $(ASSETS) $(STATIK)

.PHONY: test clean