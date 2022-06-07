windows_amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o packrat_amd64.exe -v ./cmd/*.go

windows_arm64:
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o packrat_arm64.exe -v ./cmd/*.go

mac_arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o packrat_darwin_arm64 -v ./cmd/*.go

mac_amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o packrat_darwin_amd64 -v ./cmd/*.go

linux_arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o packrat_linux_arm64 -v ./cmd/*.go

linux_amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o packrat_linux_amd64 -v ./cmd/*.go

all:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o packrat_amd64.exe -v ./cmd/*.go
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o packrat_arm64.exe -v ./cmd/*.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o packrat_darwin_arm64 -v ./cmd/*.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o packrat_darwin_amd64 -v ./cmd/*.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o packrat_linux_arm64 -v ./cmd/*.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o packrat_linux_amd64 -v ./cmd/*.go

.PHONY:test
test:
	go test -race -v ./...