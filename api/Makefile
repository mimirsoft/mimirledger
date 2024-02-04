
GOPACKAGES := $(shell go list ./...| egrep -v '/vendor')

TEST_MODULES ?= $(GOPACKAGES)
# ----------------
# INTERNAL TARGETS
#
lint-direct: # [INTERNAL]
	@golangci-lint run

start-web: # [INTERNAL]
	air -c .air.toml

test-direct: # [INTERNAL]
	go test -bench -v -race $(TEST_MODULES) $(TESTRUN)

