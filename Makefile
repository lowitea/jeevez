.DEFAULT_GOAL := help

.PHONY: help
help:  ## Список команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

OVERRIDE=`test -f deploy/docker-compose.override.yml && \
          echo '-f deploy/docker-compose.override.yml'`

.PHONY: lint
lint:  ## Запуск линтинга
	golangci-lint run

.PHONY: stop-testdb
stop-testdb:  ## Остановка тестовой базы данных
	docker rm -f psql-jeevez-testdb || true

.PHONY: start-testdb
start-testdb: stop-testdb  ## Запуск тестовой базы данных
	docker run \
		--name psql-jeevez-testdb \
		--rm \
		--network host \
		-d \
		-e POSTGRES_PASSWORD=test \
		-e POSTGRES_USER=test \
		-e POSTGRES_DB=jeevez_test \
		postgres:13

.PHONY: tests
tests: start-testdb
	sleep 1
	go test -cover -count=1 -p 1 -race ./...
	make stop-testdb

.PHONY: check
check: lint tests ## Запуск проверок проекта

.PHONY: build
build:  ## Сборка приложения
	GOOS=linux go build -o jeevez ./cmd/jeevez/main.go

.PHONY: run
run:  ## Запуск проекта
	go run cmd/jeevez/main.go

start-devdb:  ## Запуск локальной базы данных для разработки
	docker run \
		--name psql-jeevez \
		--rm \
		--network host \
		-d \
		-e POSTGRES_PASSWORD=jeevez \
		-e POSTGRES_USER=jeevez \
		-e POSTGRES_DB=jeevez \
		-v jeevez_db:/var/lib/postgresql/data \
		postgres:13

host-setup:  ## Настраивает хост для выкатки
	ansible-playbook \
		-i deploy/host/ansible/hosts.yaml \
		deploy/host/ansible/plays/base.yaml \
		deploy/host/ansible/plays/docker.yaml \
		deploy/host/ansible/plays/user.yaml
