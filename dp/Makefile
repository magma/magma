SHELL := /bin/bash
export MINIKUBE_HOME ?= $(shell echo ~/.minikube)
export KUBECONFIG ?= $(shell echo ~/.kube/config)
export CERTS := $(shell mktemp -d /tmp/certs.XXXXXXXXXX)
export DP_MODULES:= active_mode_controller configuration_controller db_service radio_controller  test_runner
export PROTO_LIST = orc8r_protos lte_protos feg_protos cwf_protos dp_protos
export MAGMA_ROOT ?= $(shell cd .. && pwd)
export PYTHON_BUILD ?= $(MAGMA_ROOT)/build
export BIN := $(PYTHON_BUILD)/bin
export SRC := $(MAGMA_ROOT)
export KUBERNETES_VERSION := 1.20.7


dev: export DATE := $(shell date +%F-%T)

ifdef REMOTE_NAMESPACE
.DEFAULT_GOAL=remote
else
.DEFAULT_GOAL=local
endif

ifdef DP_METRICS
APPEND_PROFILES=,metrics
endif

ifdef DP_CERTIFICATION
APPEND_PROFILES:=$(APPEND_PROFILES),crl,harness
endif

.PHONY: local
local: init dev

.PHONY: integration_tests
integration_tests: orc8r_integration_tests

.PHONY: orc8r_integration_tests
orc8r_integration_tests: init_orc8r _ci_init cleanup_test_db
	kubectl delete jobs -l jobName=domain-proxy-db-setup --ignore-not-found=true
	kubectl delete pods test-runner-orc8r --ignore-not-found=true
	skaffold dev --tail=false -p orc8r-deployment,integration-tests$(APPEND_PROFILES) --trigger=manual

.PHONY: dev_integration_tests
dev_integration_tests: init_orc8r _ci_init cleanup_test_db
	kubectl delete jobs -l jobName=domain-proxy-db-setup --ignore-not-found=true
	kubectl delete pods test-runner-dev --ignore-not-found=true
	skaffold dev --tail=false -p orc8r-deployment,integration-tests,integration-tests-dev$(APPEND_PROFILES) --trigger=manual

.PHONY: orc8r
orc8r: init_orc8r dev_orc8r

.PHONY: remote
remote: _remote_init dev

.PHONY: init
init: start_minikube _ci_init

.PHONY: init_orc8r
init_orc8r: start_minikube_orc8r _ci_init

.PHONY: _remote_init
_remote_init: _ci_init
ifneq ($(and $(DP_DOCKERHUB_TOKEN),$(DP_DOCKERHUB_USERNAME)),)
	@echo "$$DP_DOCKERHUB_TOKEN" | docker login --username $$DP_DOCKERHUB_USERNAME --password-stdin
	kubectx dp_orc8r
export SKAFFOLD_DEFAULT_REPO=docker.io/$(DP_DOCKERHUB_USERNAME)
else
	$(error DP_DOCKERHUB_TOKEN and DP_DOCKERHUB_USERNAME are not set)
endif

.PHONY: start_minikube
start_minikube:
	minikube start --addons=metrics-server --kubernetes-version=$(KUBERNETES_VERSION)

.PHONY: start_minikube_orc8r
start_minikube_orc8r:
	minikube start --addons=metrics-server --memory='8g' --cpus=6 --kubernetes-version=$(KUBERNETES_VERSION)
	minikube ssh sudo ip link set docker0 promisc on

.PHONY: _switch_k8s_context
_switch_k8s_context:
	kubectx minikube

.PHONY: clean
clean:
	minikube delete
	rm -rf tools/deployment/certificates/certs

.PHONY: dev_orc8r
dev_orc8r:
	skaffold dev -p orc8r-deployment,nms-deployment$(APPEND_PROFILES) --trigger=manual

.PHONY: dev
dev:
ifdef CI
ifdef REMOTE_NAMESPACE
	#MYVAR=$(DATE) skaffold dev --force=true
	skaffold dev --force=true -p remote -p integration-tests
else
	skaffold dev --force=true
endif
else
	skaffold dev --trigger=manual
endif

.PHONY: _build_ci
_build_ci: _install_skaffold_ci
	# Disable skaffold telemetry
	skaffold config set --global collect-metrics false
ifdef SERVICE
	skaffold build -b $(SERVICE)
else
	skaffold build
endif

.PHONY: _install_skaffold_ci
_install_skaffold_ci:
# Check whether skaffold is present on PATH
ifeq (, $(shell which skaffold))
	curl -Lo /tmp/skaffold https://storage.googleapis.com/skaffold/releases/v1.39.2/skaffold-linux-amd64 && \
        sudo install /tmp/skaffold /usr/local/bin/
endif

.PHONY: _ci_init
_ci_init: _generate_certificates _generate_harness_config
	kubectl delete secret certificates --ignore-not-found
	kubectl create secret generic certificates --from-file=tools/deployment/certificates/certs
	kubectl delete \
		service/orc8r-fluentd-forward \
		deployment.apps/orc8r-fluentd-forward \
		configmap/orc8r-fluentd-forward-configs \
		deployment/fake-sas-deployment \
		--ignore-not-found
	helm uninstall domain-proxy orc8r orc8r-lte || true
ifdef DP_CERTIFICATION
	kubectl delete secret private --ignore-not-found
	kubectl create secret generic private --from-file=tools/deployment/certificates/certs/private
	kubectl delete secret crl --ignore-not-found
	kubectl create secret generic crl --from-file=tools/deployment/certificates/certs/crl
	kubectl delete secret dbcbsd-ca --ignore-not-found
	kubectl create secret generic dbcbsd-ca --from-file=tools/deployment/certificates/certs/db/cbsd_ca
	kubectl delete secret dbproxy-ca --ignore-not-found
	kubectl create secret generic dbproxy-ca --from-file=tools/deployment/certificates/certs/db/proxy_ca
	kubectl delete secret dbrevoked-cbsd-ca --ignore-not-found
	kubectl create secret generic dbrevoked-cbsd-ca --from-file=tools/deployment/certificates/certs/db/revoked_cbsd_ca
	kubectl delete secret dbrevoked-proxy-ca --ignore-not-found
	kubectl create secret generic dbrevoked-proxy-ca --from-file=tools/deployment/certificates/certs/db/revoked_proxy_ca
	kubectl delete secret dbrevoked-sas-ca --ignore-not-found
	kubectl create secret generic dbrevoked-sas-ca --from-file=tools/deployment/certificates/certs/db/revoked_sas_ca
	kubectl delete secret dbroot-ca --ignore-not-found
	kubectl create secret generic dbroot-ca --from-file=tools/deployment/certificates/certs/db/root_ca
	kubectl delete secret dbsas-ca --ignore-not-found
	kubectl create secret generic dbsas-ca --from-file=tools/deployment/certificates/certs/db/sas_ca
endif
	#kubectl apply -f ./tools/deployment/vendor/postgresql.yml
	#kubectl rollout status --watch --timeout=600s statefulset/postgres-database


.PHONY: _generate_certificates
_generate_certificates:
	tools/deployment/certificates/generate_fake_certs.sh
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/device_c.cert \
	cloud/helm/dp/charts/domain-proxy/certificates/configuration_controller/device_c.cert
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/device_c.key \
	cloud/helm/dp/charts/domain-proxy/certificates/configuration_controller/device_c.key
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/ca.cert \
	cloud/helm/dp/charts/domain-proxy/certificates/configuration_controller/ca.cert

.PHONY: _generate_ci_certificates
_generate_ci_certificates: _generate_certificates
	cp -R $(CURDIR)/tools/deployment/certificates/certs/* $(CERTS)
	sudo chown -R 65534:65534 $(CERTS)

.PHONY: _generate_harness_config
_generate_harness_config:
	@set -e; \
	kubectl delete configmap harness-config --ignore-not-found; \
	kubectl create configmap harness-config \
	--from-file=tools/deployment/vendor/sas.cfg;

.PHONY: _ci_integration_tests_orc8r
_ci_integration_tests_orc8r: _install_skaffold_ci
	skaffold run -p orc8r-deployment,integration-tests$(APPEND_PROFILES)
	kubectl logs test-runner-orc8r --timestamps=true -f | tee /tmp/integration-tests-results/$@.txt
	@set -e;\
	sleep 2;\
	SUCCESS=$$(kubectl get pods test-runner-orc8r -o jsonpath='{.status.phase}');\
	if [[ $$SUCCESS != 'Succeeded' ]]; then kubectl get pods test-runner-orc8r -o yaml && exit 1; fi

.PHONY: _setup_db
_setup_db: _postgres_db_setup

.PHONY: _postgres_db_setup
_postgres_db_setup: _switch_k8s_context
	kubectl apply -f ./tools/deployment/vendor/postgresql.yml

.PHONY: _delete_db
_delete_db: _postgres_db_delete

.PHONY: _postgres_db_delete
_postgres_db_delete: _switch_k8s_context
	kubectl delete -f ./tools/deployment/vendor/postgresql.yml
	kubectl delete pvc -l app=postgres-database

.PHONY: _ci_chart_smoke_tests
_ci_chart_smoke_tests: _install_skaffold_ci
	skaffold run
	sleep 5
	helm test --timeout 10m domain-proxy

.PHONY: migration
migration:
	docker rm -f db
	@docker run -d -p 5532:5432 --name db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dp postgres \
	2> /dev/null || true
	@sleep 5
	export SQLALCHEMY_DB_URI=postgresql+psycopg2://postgres:postgres@localhost:5532/dp; \
	export PYTHONPATH=$(SRC)/dp/cloud/python; \
	cd cloud/python/magma/db_service/migrations; \
	alembic -c ./alembic.ini upgrade head; \
	alembic -c ./alembic.ini revision --autogenerate
	docker rm -f db

.PHONY: build_db
build_db:
	docker rm -f db
	@docker run -d -p 5532:5432 --name db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dp postgres \
	2> /dev/null || true
	@sleep 5
	export SQLALCHEMY_DB_URI=postgresql+psycopg2://postgres:postgres@localhost:5532/dp; \
	export PYTHONPATH=$(SRC)/dp/cloud/python; \
	cd cloud/python/magma/db_service/migrations; \
	alembic -c ./alembic.ini upgrade head; \
	cd ..; \
	python3 db_initialize.py; \

.PHONY: cleanup_test_db
cleanup_test_db: _switch_k8s_context
	kubectl exec -it postgres-database-0 -- \
		psql -U postgres -c \
		"drop database dp_test with (force);" || true ;\
	kubectl exec -it postgres-database-0 -- \
		psql -U postgres -c \
		"create database dp_test;" || true;

LOGS_COUNT ?= 1000
CBSDS_COUNT ?= 1000
LOG_FILE ?= performance_tests.log
.PHONY: prepare_performance_test_data
prepare_performance_test_data:
	@docker run -d -p 5532:5432 --name db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dp postgres \
	2> /dev/null || true
	@sleep 5
	export LOG_FILE=$(LOG_FILE); \
	export SQLALCHEMY_DB_URI=postgresql+psycopg2://postgres:postgres@localhost:5532/dp; \
	export PYTHONPATH=$(SRC)/dp/cloud/python; \
	cd tools/performance_tests; \
	python3 performance_test_setup.py -c=$(CBSDS_COUNT) -l=$(LOGS_COUNT); \

.PHONY: setup_performance_tests
setup_performance_tests: build_db prepare_performance_test_data

LIMIT ?= 100
OFFSET ?= 0
.PHONY: run_performance_tests
run_performance_tests:
	export LOG_FILE=$(LOG_FILE); \
	export DB_HOST=localhost; \
	export DB_USER=postgres; \
	export DB_PASSWORD=postgres; \
	export DB_PORT=5532; \
	export DB_NAME=dp; \
	export PYTHONPATH=$(SRC)/dp/cloud/python; \
	cd tools/performance_tests; \
	python3 run_performance_tests.py -l=$(LIMIT) -o=$(OFFSET); \

.PHONY: protos
protos:: $(PROTO_LIST) prometheus_proto
	@find $(PYTHON_BUILD)/gen -type d | tail -n +2 | sed '/__pycache__/d' | xargs -I % touch "%/__init__.py"
$(PROTO_LIST): %_protos:
	@echo "Generating python code for $* .proto files"
	@mkdir -p $(PYTHON_BUILD)/gen
	python $(SRC)/protos/gen_protos.py $(SRC)/$*/protos/ $(MAGMA_ROOT),$(MAGMA_ROOT)/orc8r/protos/prometheus $(SRC) $(PYTHON_BUILD)/gen/

.PHONY: prometheus_proto
prometheus_proto:
	python $(SRC)/protos/gen_prometheus_proto.py $(MAGMA_ROOT) $(PYTHON_BUILD)/gen

.PHONY: connect
connect:
	docker build -t ssh tools/ssh_tunnel
	docker run -d -e REMOTE_PORT=$(REMOTE_PORT) -e JUMPHOST=$(JUMPHOST) -it \
	--rm --net=host -v $(KEY):/root/.ssh/id_rsa --name ssh ssh
