DOCKER_SERVER=task-server

build:
	docker-compose build ${DOCKER_SERVER}

swag:
	swag init -g cmd/server/main.go

test:
	go test -v ./...

run-all:
	docker-compose up

run:
	docker-compose up ${DOCKER_SERVER}

create-user:
	docker exec -it ${DOCKER_SERVER} ./createuser -username $(username) -password $(password)