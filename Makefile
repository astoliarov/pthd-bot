build:
	go build -v -o bin/bot ./cmd/bot

run/bot:
	./bin/bot

docker/build:
	docker buildx build -t bghji/teamkillbot . --platform=linux/amd64

docker/push:
	docker push bghji/teamkillbot

test:
	go test ./...
