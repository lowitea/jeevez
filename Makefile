.PHONY:
	help
	check
	build
	run
	deploy

.DEFAULT_GOAL := help

help:  ## Список команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

OVERRIDE=`test -f deploy/docker-compose.override.yml && \
          echo '-f deploy/docker-compose.override.yml'`

check:  ## Запуск проверок проекта
	golangci-lint run

build:  ## Сборка приложения
	CGO_ENABLED=0 \
	GOOS=linux \
	go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd/jeevez/main.go

run:  ## Запуск проекта
	go run cmd/jeevez/main.go
