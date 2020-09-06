GOPATH?=$(HOME)/go

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/sadwave-events-tg ./cmd/main.go

docker-build:
	docker build -t punxlab/sadwave-events-tg .

docker-run:
	docker run punxlab/sadwave-events-tg:latest

docker-push:
	docker push punxlab/sadwave-events-tg:latest