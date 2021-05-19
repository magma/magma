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
	$(MAGMA_ROOT)/wifi \
	$(MAGMA_ROOT)/fbinternal
$(warning MAGMA_MODULES is undefined, using default: $(MAGMA_MODULES))
endif

SWAGGER_ROOT := $(MAGMA_ROOT)/orc8r/cloud/go/obsidian/swagger
SWAGGER_COMMON := swagger-common.yml
SWAGGER_V1_ROOT := $(SWAGGER_ROOT)/v1
SWAGGER_V1_YML := $(SWAGGER_V1_ROOT)/swagger.yml
SWAGGER_V1_CLIENT := $(SWAGGER_V1_ROOT)/client
SWAGGER_V1_CONFIG := $(MAGMA_ROOT)/orc8r/cloud/docker/controller/goswagger-config.yml
SWAGGER_V1_SPECS_DIR := $(MAGMA_ROOT)/orc8r/cloud/swagger/specs
SWAGGER_V1_PARTIAL_SPECS_DIR := $(SWAGGER_V1_SPECS_DIR)/partial
SWAGGER_V1_COMMON_DIR := $(SWAGGER_V1_SPECS_DIR)/common
SWAGGER_V1_STANDALONE_DIR := $(SWAGGER_V1_SPECS_DIR)/standalone

export SWAGGER_ROOT
export SWAGGER_COMMON
export SWAGGER_V1_ROOT
export SWAGGER_V1_CLIENT
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
	GEN_LIST += $(MODULE)_gen
	LINT_LIST += $(MODULE)_lint
	PLUGIN_LIST += $(MODULE)_plugin
	SWAGGER_LIST += $(MODULE)_swagger
	TEST_LIST += $(MODULE)_test
	TIDY_LIST += $(MODULE)_tidy
	TOOLS_LIST += $(MODULE)_tools

	MIGRATION_PLUGIN_LIST += $(MODULE)_migrationplugin
endef
_ := $(foreach module, $(MAGMA_MODULES), $(eval $(call add_module,$(module))))

###################
## General rules ##
###################

precommit: fmt lint test

build: $(BUILD_LIST) migration_plugin swagger
$(BUILD_LIST): %_build:
	make -C $*/cloud/go build

clean: $(CLEAN_LIST)
$(CLEAN_LIST): %_clean:
	make -C $*/cloud/go clean

clean_gen: clean_gen_pre $(CLEAN_GEN_LIST)
clean_gen_pre:
	rm -rf $(SWAGGER_V1_CLIENT)
	mkdir -p $(SWAGGER_V1_CLIENT)

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
	if [ -d $*/lib/go ] ; then make -C $*/lib/go fmt ; fi

gen: tools $(GEN_LIST)
$(GEN_LIST): %_gen:
	make -C $*/cloud/go gen

fullgen: clean_gen gen swagger tidy

lint: lint_tools $(LINT_LIST)
$(LINT_LIST): %_lint:
	make -C $*/cloud/go lint
	if [ -d $*/lib/go ] ; then make -C $*/lib/go lint ; fi

test: tools $(TEST_LIST)
$(TEST_LIST): %_test:
	make -C $*/cloud/go test
	if [ -d $*/lib/go ] ; then make -C $*/lib/go test ; fi

tidy: $(TIDY_LIST)
$(TIDY_LIST): %_tidy:
	make -C $*/cloud/go tidy
	if [ -d $*/lib/go ] ; then make -C $*/lib/go tidy ; fi

tools: $(TOOLS_LIST)
$(TOOLS_LIST): %_tools:
	make -C $*/cloud/go tools

lint_tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $$(go env GOPATH)/bin v1.35.2

######################
## Swagger/API docs ##
######################

swagger: swagger_directories $(SWAGGER_LIST)
	cp $(MAGMA_ROOT)/orc8r/cloud/go/models/swagger-common.yml $(SWAGGER_V1_COMMON_DIR)/$(SWAGGER_COMMON)
	combine_swagger --in=$(SWAGGER_V1_PARTIAL_SPECS_DIR) --common=$(SWAGGER_V1_COMMON_DIR)/$(SWAGGER_COMMON) --out=$(SWAGGER_V1_YML)
	swagger validate $(SWAGGER_V1_YML)

swagger_directories:
	rm -rf $(SWAGGER_V1_SPECS_DIR)
	mkdir -p $(SWAGGER_V1_COMMON_DIR) $(SWAGGER_V1_PARTIAL_SPECS_DIR) $(SWAGGER_V1_STANDALONE_DIR)


$(SWAGGER_LIST): %_swagger:
	make -C $*/cloud/go swagger_tools copy_swagger_files

##########
## Misc ##
##########

list:  # List all commands
	@echo -e "\nAvailable commands:\n"
	@grep '^[[:lower:]].*:.*#.*' Makefile | awk -F ':|#' '{ printf "\t%-30s%s\n", $$1, $$3 }'
	@echo -e "\n"

migration_plugin: $(MIGRATION_PLUGIN_LIST)
$(MIGRATION_PLUGIN_LIST): %_migrationplugin:
	make -C $*/cloud/go migration_plugin
