SHELL := /bin/bash

#########################
## Essential variables ##
#########################

ifndef MAGMA_ROOT
MAGMA_ROOT := /home/$(USER)/magma
$(warning MAGMA_ROOT is undefined, using default: $(MAGMA_ROOT))
endif
export MAGMA_ROOT

ifndef MAGMA_MODULES
MAGMA_MODULES := \
	$(MAGMA_ROOT)/orc8r \
	$(MAGMA_ROOT)/lte \
	$(MAGMA_ROOT)/feg \
	$(MAGMA_ROOT)/cwf \
	$(MAGMA_ROOT)/dp
$(warning MAGMA_MODULES is undefined, using default: $(MAGMA_MODULES))
endif

# To override this on macOS with Homebrew-installed protoc, set something like
# PROTO_INCLUDES=/usr/local/Cellar/protobuf/3.10.0/include/google/
ifndef PROTO_INCLUDES
PROTO_INCLUDES := /usr/include
endif
export PROTO_INCLUDES

SWAGGER_ROOT := $(MAGMA_ROOT)/orc8r/cloud/go/obsidian/swagger
SWAGGER_COMMON := swagger-common.yml
SWAGGER_V1_ROOT := $(SWAGGER_ROOT)/v1
SWAGGER_V1_YML := $(SWAGGER_V1_ROOT)/swagger.yml
SWAGGER_V1_CONFIG := $(MAGMA_ROOT)/orc8r/cloud/docker/controller/goswagger-config.yml
SWAGGER_V1_SPECS_DIR := $(MAGMA_ROOT)/orc8r/cloud/swagger/specs
SWAGGER_V1_PARTIAL_SPECS_DIR := $(SWAGGER_V1_SPECS_DIR)/partial
SWAGGER_V1_COMMON_DIR := $(SWAGGER_V1_SPECS_DIR)/common
SWAGGER_V1_STANDALONE_DIR := $(SWAGGER_V1_SPECS_DIR)/standalone
SWAGGER_NMS_OUT := $(MAGMA_ROOT)/nms/generated/MagmaAPIBindings.js
SWAGGER_NMS_OUT_TS := $(MAGMA_ROOT)/nms/generated-ts/

OS_NAME := $(shell uname)

export SWAGGER_ROOT
export SWAGGER_COMMON
export SWAGGER_V1_ROOT
export SWAGGER_V1_CONFIG
export SWAGGER_V1_SPECS_DIR
export SWAGGER_V1_COMMON_DIR
export SWAGGER_V1_PARTIAL_SPECS_DIR
export SWAGGER_V1_STANDALONE_DIR

COVER_DIR := $(MAGMA_ROOT)/orc8r/cloud/coverage
export COVER_DIR

#############
## Modules ##
#############

define add_module
	$(eval MODULE := $1)

	BUILD_LIST += $(MODULE)_build
	CLEAN_LIST += $(MODULE)_clean
	CLEAN_GEN_LIST += $(MODULE)_cleangen
	COVER_LIST += $(MODULE)_cover
	DOWNLOAD_LIST += $(MODULE)_download
	FMT_LIST += $(MODULE)_fmt
	GOIMPORTS_LIST += $(MODULE)_goimports
	GEN_LIST += $(MODULE)_gen
	GEN_PROTO_LIST += $(MODULE)_gen_proto
	LINT_LIST += $(MODULE)_lint
	PLUGIN_LIST += $(MODULE)_plugin
	SWAGGER_LIST += $(MODULE)_swagger
	TEST_LIST += $(MODULE)_test
	TIDY_LIST += $(MODULE)_tidy
	TOOLS_LIST += $(MODULE)_tools
endef
_ := $(foreach module, $(MAGMA_MODULES), $(eval $(call add_module,$(module))))

###################
## General rules ##
###################

precommit: fmt lint test

build: $(BUILD_LIST) swagger
$(BUILD_LIST): %_build:
	make -C $*/cloud/go build

clean: $(CLEAN_LIST)
$(CLEAN_LIST): %_clean:
	make -C $*/cloud/go clean

clean_gen: $(CLEAN_GEN_LIST)
$(CLEAN_GEN_LIST): %_cleangen:
	make -C $*/cloud/go clean_gen

cover: cover_pre $(COVER_LIST) cover_post
$(COVER_LIST): %_cover:
	make -C $*/cloud/go cover
cover_pre:
	mkdir -p $(COVER_DIR)
cover_post:
	gocovmerge $(COVER_DIR)/** > $(COVER_DIR)/all.gocov

download: $(DOWNLOAD_LIST)
$(DOWNLOAD_LIST): %_download:
	make -C $*/cloud/go download

fmt: $(FMT_LIST)
$(FMT_LIST): %_fmt:
	make -C $*/cloud/go fmt

fullgen: clean_gen gen swagger nms_fullgen tidy

gen: tools order_imports gen_protos $(GEN_LIST)
$(GEN_LIST): %_gen:
	make -C $*/cloud/go gen

order_imports: $(GOIMPORTS_LIST)
$(GOIMPORTS_LIST): %_goimports:
	make -C $*/cloud/go order_imports

gen_protos: $(GEN_PROTO_LIST)
$(GEN_PROTO_LIST): %_gen_proto:
	make -C $*/cloud/go gen_protos

lint: tools_lint $(LINT_LIST)
$(LINT_LIST): %_lint:
	make -C $*/cloud/go lint

test: tools $(TEST_LIST)
$(TEST_LIST): %_test:
	make -C $*/cloud/go test

tidy: $(TIDY_LIST)
$(TIDY_LIST): %_tidy:
	make -C $*/cloud/go tidy

tools: $(TOOLS_LIST)
$(TOOLS_LIST): %_tools:
	make -C $*/cloud/go tools

tools_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $$(go env GOPATH)/bin v1.45.0

######################
## Swagger/API docs ##
######################

swagger: swagger_directories $(SWAGGER_LIST)
	cp $(MAGMA_ROOT)/orc8r/cloud/go/models/swagger-common.yml $(SWAGGER_V1_COMMON_DIR)/$(SWAGGER_COMMON)
	swaggergen --combine --in=$(SWAGGER_V1_PARTIAL_SPECS_DIR) --common=$(SWAGGER_V1_COMMON_DIR)/$(SWAGGER_COMMON) --out=$(SWAGGER_V1_YML)
	swagger validate $(SWAGGER_V1_YML)

swagger_directories:
	rm -rf $(SWAGGER_V1_SPECS_DIR)
	mkdir -p $(SWAGGER_V1_COMMON_DIR) $(SWAGGER_V1_PARTIAL_SPECS_DIR) $(SWAGGER_V1_STANDALONE_DIR)

$(SWAGGER_LIST): %_swagger:
	make -C $*/cloud/go swagger_tools copy_swagger_files

#############
## NMS gen ##
#############

nms_fullgen: nms_prereqs nms_gen

nms_prereqs:
ifeq ($(OS_NAME),Linux)
	$(MAKE) nms_prereqs_ubuntu
else ifeq ($(OS_NAME),Darwin)
	$(MAKE) nms_prereqs_osx
endif

nms_prereqs_ubuntu: nms_node_ubuntu
	NODE="$(shell which nodejs)"; ln -s $$NODE /usr/local/bin/node
	# update nodejs to latest
	curl -sL https://deb.nodesource.com/setup_lts.x | bash -
	apt-get install -y nodejs
	# get codegen dependency
	npm install --global yarn
	yarn add swagger-js-codegen
	yarn add @openapitools/openapi-generator-cli

nms_node_ubuntu:
	# install pre-reqs
	apt update
	# Handle certificate verification
	apt-get install -y libgnutls30
	# Get a newer version of nodejs
	curl -fSL https://deb.nodesource.com/setup_14.x | bash -
	apt-get install -y nodejs
	# And now install npm
	apt install -y aptitude
	aptitude install -y npm
	# symlink nodejs executable
	mkdir -p /usr/local/bin

nms_prereqs_osx:
	node --version || brew install node
	npm --version || brew install npm
	yarn --version || npm install --global yarn
	yarn add swagger-js-codegen
	yarn add @openapitools/openapi-generator-cli

nms_gen:
	# generate Swagger API bindings for NMS
	$(MAGMA_ROOT)/nms/scripts/generateAPIFromSwagger.sh $(SWAGGER_V1_YML) $(SWAGGER_NMS_OUT)
	$(MAGMA_ROOT)/nms/scripts/generateAPIFromSwaggerForTypeScript.sh -f -i $(SWAGGER_V1_YML) -o $(SWAGGER_NMS_OUT_TS)
