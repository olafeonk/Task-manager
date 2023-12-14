DOCKER_SERVER=task-server

build:
	docker-compose build ${DOCKER_SERVER}

swag:
	swag init -g cmd/server/main.go

run-all:
	docker-compose up

run:
	docker-compose up ${DOCKER_SERVER}

