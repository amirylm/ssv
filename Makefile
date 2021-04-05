ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

ifndef $(BUILD_PATH)
    BUILD_PATH="/go/bin/ssvnode"
    export BUILD_PATH
endif

ifneq (,$(wildcard ./.env))
    include .env
endif

#Lint
.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint:
	./bin/golangci-lint run -v ./...

#Test

#Test
.PHONY: full-test
full-test:
	@echo "Running the full test..."
	@go test -tags blst_enabled -timeout 20m -cover -race -p 1 -v ./...



# TODO: Intgrate use of short flag (unit tests) + running tests through docker
#.PHONY: unittest
#unittest:
#	@go test -v -short -race ./...

#.PHONY: full-test-local
#full-test-local:
#	@docker-compose -f test.docker-compose.yaml up -d mysql_test
#	@make full-test
#	# @docker-compose -f test.docker-compose.yaml down --volumes
#
#
#.PHONY: docker-test
#docker-test:
#	@docker-compose -f test.docker-compose.yaml up -d mysql_test
#	@docker-compose -f test.docker-compose.yaml up --build --abort-on-container-exit
#	@docker-compose -f test.docker-compose.yaml down --volumes


.PHONY: start-node
start-node:
	@echo "Build ${BUILD_PATH}"
ifdef DEBUG_PORT
	@echo "Running node-${NODE_ID} in debug mode"
	@dlv  --continue --accept-multiclient --headless --listen=:${DEBUG_PORT} --api-version=2 exec \
	 ${BUILD_PATH} start-node -- --node-id=${NODE_ID} --private-key=${SSV_PRIVATE_KEY} --validator-key=${VALIDATOR_PUBLIC_KEY} --beacon-node-addr=${BEACON_NODE_ADDR} --network=${NETWORK} --discovery-type=mdns --val=${CONSENSUS_TYPE}

else
	@echo "Running node (${NODE_ID})"
	${BUILD_PATH} start-node --node-id=${NODE_ID} --private-key=${SSV_PRIVATE_KEY} --validator-key=${VALIDATOR_PUBLIC_KEY} --beacon-node-addr=${BEACON_NODE_ADDR} --network=${NETWORK} --val=${CONSENSUS_TYPE} --host-dns=${HOST_DNS}
endif

NODES=ssv-node-1 ssv-node-2 ssv-node-3 ssv-node-4
.PHONY: docker
docker:
	@echo "nodes $(NODES)"
	@docker-compose up --build $(NODES)

DEBUG_NODES=ssv-node-1-dev ssv-node-2-dev ssv-node-3-dev ssv-node-4-dev
.PHONY: docker-debug
docker-debug:
	@docker-compose up --build $(DEBUG_NODES)

.PHONY: stop
stop:
	@docker-compose  down


.PHONY: start-boot-node
start-boot-node:
	@echo "Running start-boot-node"
	${BUILD_PATH} start-boot-node --private-key=${BOOT_NODE_PRIVATE_KEY}
