# REQUIRED SECTION
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
include $(ROOT_DIR)/.mk-lib/common.mk
# END OF REQUIRED SECTION

.PHONY: help dependencies up start stop restart logs status ps clean test pytest

dependencies: check-dependencies ## Check dependencies

build: ## Build all or c=<name> containers
	@echo "DEPRECATED: run the command 'docker-compose build' instead"

up: ## Start all or c=<name> containers in foreground
	@echo "DEPRECATED: run the command 'docker-compose up' instead"

start: ## Start all or c=<name> containers in background
	@echo "DEPRECATED: run the command 'docker-compose up' instead"

stop: ## Stop all or c=<name> containers
	@echo "DEPRECATED: run the command 'docker-compose stop' instead"

restart: ## Restart all containers
	@echo "DEPRECATED: run the command 'docker-compose build && docker-compose up' instead"

logs: ## Show all or c=<name> container logs
	@echo "DEPRECATED: run the command 'docker-compose logs' instead"

status: ## Show status of containers
	@echo "DEPRECATED: run the command 'docker-compose ps' instead"

ps: ## Alias of status
	@echo "DEPRECATED: run the command 'docker-compose ps' instead"

clean: ## Clean all data
	@echo "DEPRECATED: run the command 'docker-compose down --remove-orphans' instead"

test: ## Run integration tests
	@echo "DEPRECATED: run the command 'docker-compose build && docker-compose up && docker-compose run --use-aliases test go test -v ./...' instead"

pytest: ## Run integration tests
	@echo "DEPRECATED: run the command 'docker-compose -f docker-compose.yaml -f docker-compose.testing.yaml build && docker-compose -f docker-compose.yaml -f docker-compose.testing.yaml up -d && docker-compose -f docker-compose.yaml -f docker-compose.testing.yaml run --use-aliases pytest' instead"

cover: ## Run go test coverage
	@-go test ../... -coverprofile ./cover.tmp >/dev/null
	@go tool cover -func=./cover.tmp | tail -n 1
	@find .. -name '*.go' | xargs wc -l | tail -n 1
	@rm ./cover.tmp


### These targets set up connecting the nms to the local docker instance of orc8r cloud controller
.PHONY: up-local copy-local-certs
MAGMA_CERT_DIR=../../../magma/.cache/test_certs
NMS_CERT_DIR=~/magma_certs

up-local: ## Point NMS at local docker orc8r container(s)
	@if test -f $(NMS_CERT_DIR)/local_api_cert -a -f $(NMS_CERT_DIR)/local_api_key; then \
	  MAGMA_API_HOST=host.docker.internal:9443 MAGMA_API_AUTH=local_ $(DOCKER_COMPOSE) up -d $(c); \
	else echo "missing orc8r cert/key in $(NMS_CERT_DIR)" && false; \
	fi

copy-local-certs: ## Copy certs from magma test_certs path to nms's magma_certs path
	@if test -f $(MAGMA_CERT_DIR)/admin_operator.pem -a \
		 -f $(MAGMA_CERT_DIR)/admin_operator.key.pem; then \
		cp $(MAGMA_CERT_DIR)/admin_operator.pem $(NMS_CERT_DIR)/local_api_cert; \
		cp $(MAGMA_CERT_DIR)/admin_operator.key.pem $(NMS_CERT_DIR)/local_api_key; \
		echo "local certs copied from $(MAGMA_CERT_DIR) to $(NMS_CERT_DIR)"; \
	else echo "missing orc8r cert/key in $(MAGMA_CERT_DIR)" && false; \
	fi

.PHONY: up-boslab-nms
BOSLAB_CERT_DIR=../../southpoll/seed_files/orc8r-cache/test_certs
BOSLAB_HOST ?= localhost
BOSLAB_PORT ?= 8002
BOSLAB_MAGMA_API_HOST ?= host.docker.internal
# curl -s -k -X GET https://localhost:8002/magma/networks -H "accept: application/json" --cert $(BOSLAB_CERT_DIR)/admin_operator.pem --key $(BOSLAB_CERT_DIR)/admin_operator.key.pem > /dev/null
up-boslab-nms: ## Point NMS at boslab docker orc8r container(s) via ssh tunnels (cd ../../../fbc/boslab && make proxynms)
	@if ! nc -z $(BOSLAB_HOST) $(BOSLAB_PORT) > /dev/null ; then \
	  echo "Orc8r connection not successful! Did you run 'make proxynms' in fbc/boslab?" && false; \
	fi

	@echo "Testing orc8r boslab certs..."
	@if curl -s -k -X GET https://$(BOSLAB_HOST):$(BOSLAB_PORT)/magma/networks -H "accept: application/json" --cert $(BOSLAB_CERT_DIR)/admin_operator.pem --key $(BOSLAB_CERT_DIR)/admin_operator.key.pem > /dev/null ; then \
	  echo "...orc8r connection validated"; \
	else \
	  echo "Orc8r connection not successful! Are certs correct? Did you run 'make proxynms' in fbc/boslab?" && false; \
	fi

	# copy boslab cert and point docker images at boslab
	cp $(BOSLAB_CERT_DIR)/admin_operator.pem $(NMS_CERT_DIR)/boslab_api_cert
	cp $(BOSLAB_CERT_DIR)/admin_operator.key.pem $(NMS_CERT_DIR)/boslab_api_key
	@if test -f $(NMS_CERT_DIR)/boslab_api_cert -a -f $(NMS_CERT_DIR)/boslab_api_key; then \
	  MAGMA_API_HOST=$(BOSLAB_MAGMA_API_HOST):$(BOSLAB_PORT) MAGMA_API_AUTH=boslab_ $(DOCKER_COMPOSE) up -d $(c); \
	else echo "missing boslab orc8r cert/key in $(BOSLAB_CERT_DIR)" && false; \
	fi

	@echo "Testing orc8r boslab certs within container..."
	@if docker-compose exec platform-server curl -s -k -X GET https://$(BOSLAB_MAGMA_API_HOST):$(BOSLAB_PORT)/magma/networks -H "accept: application/json" --cert /magma_certs/boslab_api_cert --key /magma_certs/boslab_api_key > /dev/null ; then \
	  echo "...orc8r connection within container validated"; \
	else \
	  echo "Orc8r connection within container not successful!" && false; \
	fi

.PHONY: up-stag-nms
up-stag-nms: ## Point NMS at staging orc8r instance
	@MAGMA_API_HOST=api-staging.magma.etagecom.io $(DOCKER_COMPOSE) up -d $(c)

.PHONY: up-prod-nms
up-prod-nms: ## Point NMS at prod orc8r instance
	@MAGMA_API_HOST=api.magma.etagecom.io $(DOCKER_COMPOSE) up -d $(c)
