# Define the binary name
BINARY_NAME=syncbuddy
# Define the build output path and filename
BINARY_OUTPUT=bin\$(BINARY_NAME).exe
# Define the Go package to build
GO_PACKAGE=./cmd/$(BINARY_NAME)

# Define the installation directory using Windows-style backslashes
INSTALL_DIR=C:\Users\eric_ekholm\bin

.PHONY: build install clean test

build:
	go build -o $(BINARY_OUTPUT) $(GO_PACKAGE)

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	copy $(BINARY_OUTPUT) $(INSTALL_DIR)
	@echo "Installation complete!"

clean:
	@if exist $(BINARY_OUTPUT) del $(subst /,\\,$(BINARY_OUTPUT))

test:
	go test ./... -cover