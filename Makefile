windows:
		go build -ldflags="-s -w" -o packrat.exe -v ./cmd/main.go

.PHONY:test
test:
		go test -race -v ./...