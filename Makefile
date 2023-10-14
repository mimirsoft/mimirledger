SERVICE_NAME := mimirledger
DOCKER_COMPOSE := docker-compose -f dev/dev.yml
DOCKER_COMPOSE_START := ${DOCKER_COMPOSE} up -d ${SERVICE_NAME}
DB_SERVICE_NAME := postgres
DB_DOCKER_COMPOSE := docker-compose -f db/db.yml
DB_DOCKER_COMPOSE_START := ${DB_DOCKER_COMPOSE} up -d ${DB_SERVICE_NAME}
DOCKER_COMPOSE_TEST := docker-compose -f dev/test.yml
LOG_TAIL_LENGTH=50


ifdef TEST_RUN
 TESTRUN := -run ${TEST_RUN}
endif

start:
	${DOCKER_COMPOSE_START}

start-db:
	${DB_DOCKER_COMPOSE_START}

stop:
	${DOCKER_COMPOSE} stop ${SERVICE_NAME}
	${DB_DOCKER_COMPOSE} down

stop-db:
	${DB_DOCKER_COMPOSE} stop ${DB_SERVICE_NAME}
	${DB_DOCKER_COMPOSE} down

restart: stop start

test: # run unit tests
	${DOCKER_COMPOSE_TEST} rm --force || true
	${DOCKER_COMPOSE_TEST} run test_mimirledger
	${DOCKER_COMPOSE_TEST} down

lint: # Run go lint
	${DOCKER_COMPOSE_TEST} run test_mimirledger ash -c "GOGC=50 make -e lint-direct"

update: stop-api docker-clean # rebuild image and restart service
	${DOCKER_COMPOSE} rm --force ${SERVICE_NAME}
	${DOCKER_COMPOSE} build ${SERVICE_NAME}
	$(MAKE) start

stop-api: # stop api service
	${DOCKER_COMPOSE} stop --force ${SERVICE_NAME} || true

logs:
	${DOCKER_COMPOSE} logs --tail $(LOG_TAIL_LENGTH) -f ${SERVICE_NAME}


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

api-shell:
	docker exec -it dev_mimirledger_1  /bin/ash

postgres-shell:
	docker exec -it dev_db-postgres_1  /bin/bash

psql:
	$(MAKE) -C db -e psql
