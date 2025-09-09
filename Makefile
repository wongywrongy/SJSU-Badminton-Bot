run:
	go run ./cmd/bot

test:
	go test ./...

build:
	go build -o bin/bot ./cmd/bot

docker-build:
	docker build -t sjsu-badminton-bot:latest .
