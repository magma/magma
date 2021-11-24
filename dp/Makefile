SHELL := /bin/bash
export MINIKUBE_HOME ?= $(shell echo ~/.minikube)
export KUBECONFIG ?= $(shell echo ~/.kube/config)
export CERTS := $(shell mktemp -d /tmp/certs.XXXXXXXXXX)
export DP_MODULES:= active_mode_controller configuration_controller db_service protocol_controller radio_controller  test_runner
export PROTO_LIST = orc8r_protos lte_protos feg_protos cwf_protos dp_protos
export MAGMA_ROOT ?= $(shell cd .. && pwd)
export PYTHON_BUILD ?= $(MAGMA_ROOT)/build
export BIN := $(PYTHON_BUILD)/bin
export SRC := $(MAGMA_ROOT)


dev: export DATE := $(shell date +%F-%T)

ifdef REMOTE_NAMESPACE
.DEFAULT_GOAL=remote
else
.DEFAULT_GOAL=local
endif

.PHONY: local
local: init dev

.PHONY: remote
remote: _remote_init dev

.PHONY: init
init: start_minikube _ci_init _contour_install

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
	minikube start --addons=metrics-server

.PHONY: clean
clean:
	minikube delete

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
	curl -Lo /tmp/skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
        sudo install /tmp/skaffold /usr/local/bin/
endif

.PHONY: _ci_init
_ci_init: _generate_certificates _generate_harness_config
	kubectl delete secret certificates --ignore-not-found
	kubectl create secret generic certificates --from-file=tools/deployment/certificates/certs
	#kubectl apply -f ./tools/deployment/vendor/postgresql.yml
	#kubectl rollout status --watch --timeout=600s statefulset/postgres-database

.PHONY: _contour_install
_contour_install:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm upgrade --install \
	--set envoy.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-internal"=0.0.0.0/0 \
	contour bitnami/contour --version 4.3.2
	kubectl wait --for=condition=Available --timeout=600s Deployment/contour-contour

.PHONY: _nginx_install
_nginx_install:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install \
	--set controller.extraEnvs[0].name="GODEBUG",controller.extraEnvs[0].value="x509ignoreCN=0" \
	nginx-ingress ingress-nginx/ingress-nginx --version 3.34.0
	kubectl wait --for=condition=Available --timeout=600s Deployment/nginx-ingress-ingress-nginx-controller
	#--set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-internal"=true

.PHONY: _generate_certificates
_generate_certificates:
	tools/deployment/certificates/generate_fake_certs.sh
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/domain_proxy_bundle.cert \
	cloud/helm/dp/charts/domain-proxy/certificates/protocol_controller/domain_proxy_bundle.cert
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/domain_proxy_server.key \
	cloud/helm/dp/charts/domain-proxy/certificates/protocol_controller/domain_proxy_server.key
	ln -s -f ../../../../../../../tools/deployment/certificates/certs/ca.cert \
	cloud/helm/dp/charts/domain-proxy/certificates/protocol_controller/ca.cert
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

.PHONY: _ci_integration_tests
_ci_integration_tests: _install_skaffold_ci _contour_install
	kubectl delete --ignore-not-found=true job test-runner-job
	skaffold run
	kubectl wait --for=condition=complete --timeout=10m job/test-runner-job & \
	kubectl wait --for=condition=failed --timeout=10m job/test-runner-job & \
	wait -n 1 2
	kubectl logs -l type=integration-tests
	@set -e;\
	SUCCESS=$$(kubectl get jobs test-runner-job -o jsonpath='{.status.succeeded}');\
	if [[ -z $$SUCCESS ]]; then SUCCESS=0; fi; \
	if [[ $$SUCCESS != '1' ]]; then exit 1; fi

.PHONY: _setup_db
_setup_db: _postgres_db_setup

.PHONY: _postgres_db_setup
_postgres_db_setup:
	kubectl apply -f ./tools/deployment/vendor/postgresql.yml

.PHONY: _delete_db
_delete_db: _postgres_db_delete

.PHONY: _postgres_db_delete
_postgres_db_delete:
	kubectl delete -f ./tools/deployment/vendor/postgresql.yml
	kubectl delete pvc -l app=postgres-database

.PHONY: _ci_chart_smoke_tests
_ci_chart_smoke_tests: _install_skaffold_ci _contour_install
	skaffold run
	helm test --timeout 10m domain-proxy

.PHONY: _ci_ingress_tests_contour
_ci_ingress_tests_contour: _install_skaffold_ci _contour_install _generate_ci_certificates
	sudo --preserve-env=MINIKUBE_HOME,KUBECONFIG minikube tunnel > /dev/null &
	skaffold run
	sleep 10
	cd tools/deployment/certificates/certs; \
	curl --key device_b.key --cert device_b.cert \
	--cacert ca.cert \
	--resolve "domain-proxy:443:$$(kubectl get svc \
	contour-envoy --output jsonpath='{.status.loadBalancer.ingress[0].ip}')" \
	https://domain-proxy/sas/v1/registration 

.PHONY: _ci_ingress_tests_nginx
_ci_ingress_tests_nginx: _install_skaffold_ci _nginx_install _generate_ci_certificates
	sudo --preserve-env=MINIKUBE_HOME,KUBECONFIG minikube tunnel > /dev/null &
	skaffold run --profile nginx-ingress
	sleep 10
	cd tools/deployment/certificates/certs; \
	curl --key device_b.key --cert device_b.cert \
	--cacert ca.cert \
	--resolve "domain-proxy:443:$$(kubectl get svc \
	nginx-ingress-ingress-nginx-controller --output jsonpath='{.status.loadBalancer.ingress[0].ip}')" \
	https://domain-proxy/sas/v1/registration 

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
