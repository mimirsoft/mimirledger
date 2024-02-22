
PG_DB_PORT ?= 5432
PG_DB_HOST ?= localhost
PG_DB_PASSWORD ?= mimirledger
PG_DB_USER ?= mimirledger
DEVDBNAME ?= mimirledger
TESTDBNAME ?= mimirledgertest
PG_CONTAINER := dev_postgres_1

DEV_PSQL_OPTS := -h ${PG_DB_HOST} -p ${PG_DB_PORT} -d ${DEVDBNAME} -U ${PG_DB_USER}
TEST_PSQL_OPTS := -h ${PG_DB_HOST} -p ${PG_DB_PORT} -d ${TESTDBNAME} -U ${PG_DB_USER}

create-users:
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres createuser -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres -s mimirledger || createuser -s mimirledger || true"
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres psql -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres -d postgres -c \"ALTER USER mimirledger WITH PASSWORD 'mimirledger';\""
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres psql -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres -d postgres -c \"ALTER USER mimirledger WITH PASSWORD 'mimirledger';\""
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres psql -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres -d postgres -c \"ALTER USER mimirledger WITH SUPERUSER;\""

create-devdb: create-users
	if docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} createdb -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U ${PG_DB_USER} ${DEVDBNAME}"; then \
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} psql ${DEV_PSQL_OPTS} --echo-errors < ./db/schema.sql;"; \
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} psql ${DEV_PSQL_OPTS} --echo-errors < ./db/initial_rows.sql;"; \
	else true; \
	fi

create-testdb: create-users
	if docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} createdb -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U ${PG_DB_USER} ${TESTDBNAME}"; then \
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} psql ${TEST_PSQL_OPTS} --echo-errors < ./db/schema.sql;" \
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=${PG_DB_PASSWORD} psql ${TEST_PSQL_OPTS} --echo-errors < ./db/initial_rows.sql;" \
	else true; \
	fi

drop-devdb:
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres psql -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres -d ${DEVDBNAME} -c \"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname = current_database() AND pid <> pg_backend_pid() \""
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres dropdb -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres ${DEVDBNAME} || true"

drop-testdb:
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=postgres dropdb -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U postgres ${TESTDBNAME} || true"

psql:
	docker exec -it ${PG_CONTAINER} sh -c "PGPASSWORD=mimirledger psql -h ${PG_DB_HOST} -p ${PG_DB_PORT} -U mimirledger -d mimirledger"



PGPASSWORD=postgres dropdb -h localhost -p 5432 -U postgres mimirledgertest || true
PGPASSWORD=postgres createdb -h localhost -p 5432 -U postgres mimirledgertest
PGPASSWORD=postgres psql -h localhost -p 5432 -d mimirledgertest -U postgres --echo-errors < ./db/schema.sql