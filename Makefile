.PHONY: build build-all clean run

BINARY := dotsetup
CMD := ./cmd/dotsetup

# Default: build for current platform
build:
	go build -o $(BINARY) $(CMD)

# Cross-compile for all target platforms
build-all:
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)_darwin_arm64 $(CMD)
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)_darwin_amd64 $(CMD)
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)_linux_amd64 $(CMD)
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY)_linux_arm64 $(CMD)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -rf dist/

# Install locally
install: build
	install -d $(HOME)/.local/bin
	install -m 755 $(BINARY) $(HOME)/.local/bin/$(BINARY)
