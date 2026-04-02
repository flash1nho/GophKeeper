VERSION := $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell git show -s --format=%as HEAD)
LDFLAGS := -X 'github.com/flash1nho/GophKeeper/pkg/version.buildVersion=$(VERSION)' \
           -X 'github.com/flash1nho/GophKeeper/pkg/version.buildCommit=$(COMMIT)' \
           -X 'github.com/flash1nho/GophKeeper/pkg/version.buildDate=$(BUILD_DATE)'

build:
	# gophkeeper-client
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/download/v1.0.0/gophkeeper-client-linux-amd64 app/client/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o releases/download/v1.0.0/gophkeeper-client-darwin-arm64 app/client/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/download/v1.0.0/gophkeeper-client-windows-amd64.exe app/client/main.go
