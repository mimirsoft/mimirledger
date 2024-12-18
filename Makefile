LOG_TAIL_LENGTH=50
MY_UID := $(shell id -u)
MY_GID := $(shell id -g)

SERVICE_NAME := mimirledger
DOCKER_COMPOSE := MY_UID=${MY_UID} MY_GID=${MY_GID} docker compose -f dev/dev.yml
DOCKER_COMPOSE_START := MY_UID=${MY_UID} MY_GID=${MY_GID} ${DOCKER_COMPOSE} up -d ${SERVICE_NAME}
DOCKER_COMPOSE_TEST := MY_UID=${MY_UID} MY_GID=${MY_GID} docker compose -f dev/test.yml
DB_SERVICE_NAME := postgres
DB_DOCKER_COMPOSE := docker compose -f db/db.yml
DB_DOCKER_COMPOSE_START := ${DB_DOCKER_COMPOSE} up -d ${DB_SERVICE_NAME}
WEB_SERVICE_NAME := web
WEB_DOCKER_COMPOSE := docker compose -f client/client.yml
WEB_DOCKER_COMPOSE_START := ${WEB_DOCKER_COMPOSE} up -d ${WEB_SERVICE_NAME}

ifdef TEST_RUN
 TESTRUN := -run ${TEST_RUN}
endif

print_vars:
	echo ${MY_UID}
	echo ${MY_GID}

start: start-db start-api start-web

start-api:
	${DOCKER_COMPOSE_START}

start-db:
	${DB_DOCKER_COMPOSE_START}

start-web:
	${WEB_DOCKER_COMPOSE_START}

build-web:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 ${WEB_DOCKER_COMPOSE} build ${WEB_SERVICE_NAME}

stop: stop-api stop-db stop-web

stop-api:
	${DOCKER_COMPOSE} stop ${SERVICE_NAME} || true
	${DOCKER_COMPOSE} down

stop-db:
	${DB_DOCKER_COMPOSE} stop ${DB_SERVICE_NAME}
	${DB_DOCKER_COMPOSE} down

stop-web:
	${WEB_DOCKER_COMPOSE} stop ${WEB_SERVICE_NAME}
	${WEB_DOCKER_COMPOSE} down

restart: stop start

test: drop-testdb create-testdb # run unit tests
	${DOCKER_COMPOSE_TEST} rm --force || true
	${DOCKER_COMPOSE_TEST} run test_mimirledger
	${DOCKER_COMPOSE_TEST} down

lint: # Run go lint
	${DOCKER_COMPOSE_TEST} run test_mimirledger ash -c "GOGC=50 make -e lint-direct"

lint-web: # Run eslint
	${WEB_DOCKER_COMPOSE} run ${WEB_SERVICE_NAME} ash -c "npm run lint"

lint-web-fix: # Run eslint
	${WEB_DOCKER_COMPOSE} run ${WEB_SERVICE_NAME} ash -c "npm run lint-fix"

update: stop-api docker-clean # rebuild image and restart service
	${DOCKER_COMPOSE} rm --force ${SERVICE_NAME}
	${DOCKER_COMPOSE} build ${SERVICE_NAME}
	$(MAKE) start

logs:
	${DOCKER_COMPOSE} logs --tail $(LOG_TAIL_LENGTH) -f ${SERVICE_NAME}

logs-db:
	${DB_DOCKER_COMPOSE} logs --tail $(LOG_TAIL_LENGTH) -f ${DB_SERVICE_NAME}

logs-web:
	${WEB_DOCKER_COMPOSE} logs --tail $(LOG_TAIL_LENGTH) -f ${WEB_SERVICE_NAME}

docker-clean: # clean out all containers (does NOT require a full rebuild)
	${DOCKER_COMPOSE} down || true
	${DOCKER_COMPOSE} rm --force || true
	docker rmi dev_mimirledger || true
	docker rmi $(shell docker images -f 'dangling=true' -q) 2>/dev/null || true

#
# Database utilities
#
create-devdb:
	$(MAKE) -C db -e create-devdb

drop-devdb:
	$(MAKE) -C db -e drop-devdb

create-testdb:
	$(MAKE) -C db -e create-testdb

drop-testdb:
	$(MAKE) -C db -e drop-testdb

psql:
	$(MAKE) -C db -e psql

dumpdatabase:
	$(MAKE) -C db -e dumpdatabase

loaddatabase:
	$(MAKE) -C db -e loaddatabase

#
# Container Shell tools
#
api-shell:
	docker exec -it dev_mimirledger_1  /bin/ash

web-shell:
	docker exec -it client-web-1  /bin/ash

postgres-shell:
	docker exec -it dev_postgres_1  /bin/bash
