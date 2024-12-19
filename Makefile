.PHONY: build up down

build:
	docker build . --file Dockerfile -t app:latest
	docker-compose up -d

down:
	docker-compose down
