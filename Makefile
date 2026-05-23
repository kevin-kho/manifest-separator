build:
	go build -o bin/manifest-seperator main.go

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
