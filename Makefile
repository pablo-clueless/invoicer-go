dev:
	air

build:
	go build -o invoicer ./...

test:
	go test ./...

start:
	go run main.go migrate && go run main.go
