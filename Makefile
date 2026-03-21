VERSION := $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell git show -s --format=%as HEAD)
LDFLAGS := -X 'github.com/flash1nho/GophKeeper/pkg/version.buildVersion=$(VERSION)' \
           -X 'github.com/flash1nho/GophKeeper/pkg/version.buildCommit=$(COMMIT)' \
           -X 'github.com/flash1nho/GophKeeper/pkg/version.buildDate=$(BUILD_DATE)'

build:
	# gophkeeper-server
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/server/linux/gophkeeper-server app/server/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/server/windows/gophkeeper-server.exe app/server/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/server/darwin/gophkeeper-server app/server/main.go

	# gophkeeper-client
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/client/linux/gophkeeper-client app/client/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/client/windows/gophkeeper-client.exe app/client/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/client/darwin/gophkeeper-client app/client/main.go
