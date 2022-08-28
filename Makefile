build:
	go build -v -o bin/bot ./cmd/bot
	go build -v -o bin/migrate ./cmd/migrate

run/bot:
	./bin/bot

run/migrate:
	./bin/migrate

docker/build:
	docker buildx build -t bghji/teamkillbot . --platform=linux/amd64

docker/push:
	docker push bghji/teamkillbot
