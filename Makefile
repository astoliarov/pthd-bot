build:
	go build -v -o bin/bot ./cmd/bot

run/bot:
	./bin/bot

docker/build:
	docker build . -t bghji/teamkillbot

docker/push:
	docker push bghji/teamkillbot
