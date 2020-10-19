.PHONY:
	help
	check
	build
	push
	pull
	run
	deploy

.DEFAULT_GOAL := help

help:  ## Список команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

OVERRIDE=`test -f deploy/docker-compose.override.yml && \
          echo '-f deploy/docker-compose.override.yml'`

check:  ## Запуск проверок проекта в докере
	docker-compose \
		-f deploy/docker-compose.autotests.yml \
	    ${OVERRIDE} \
		--project-directory . \
		run --rm autotests
	docker-compose  \
		-f deploy/docker-compose.autotests.yml \
		${OVERRIDE} \
		--project-directory . \
		stop

build:  ## Сборка образа докера
	test -z "${GITHUB_SHA}" || echo "${GITHUB_SHA}" > .git_commit_sha
	test -z "${GITHUB_REF}" || echo "${GITHUB_REF}" > .git_ref_tag
	docker-compose \
		-f deploy/docker-compose.yml \
	    ${OVERRIDE} \
		--project-directory . \
		build --pull --compress

push:  ## Публикация образа докера в хаб
	docker-compose -f deploy/docker-compose.yml push

pull:  ## Скачивание образа докера с хаба
	docker-compose -f deploy/docker-compose.yml pull

run:  ## Запуск проекта в докере
	docker-compose \
	    -f deploy/docker-compose.yml \
	    ${OVERRIDE} \
		--project-directory . \
	    up
