windows:
		go build -ldflags="-s -w" -o packrat.exe -v ./cmd/windows/main.go

.PHONY:test
test:
		go test -race -v ./...