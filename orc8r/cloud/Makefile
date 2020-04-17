SHELL := /bin/bash

ifndef MAGMA_MODULES
$(error MAGMA_MODULES is undefined)
endif

##########################################
## Some Useful Directory/File Variables ##
##########################################
# Allow MAGMA_ROOT and PLUGIN_DIR to be set from environment to support travis
ifndef MAGMA_ROOT
MAGMA_ROOT := /home/$(USER)/magma
endif
export MAGMA_ROOT

SWAGGER_ROOT := $(MAGMA_ROOT)/orc8r/cloud/docker/controller/apidocs
SWAGGER_COMMON := swagger-common.yml
SWAGGER_V1_ROOT := $(SWAGGER_ROOT)/v1
SWAGGER_V1_YML := $(SWAGGER_ROOT)/v1/swagger.yml
SWAGGER_V1_TEMPLATE := $(SWAGGER_ROOT)/v1/swagger-template.yml
SWAGGER_V1_TEMP_GEN := $(SWAGGER_ROOT)/v1/temp
export SWAGGER_ROOT
export SWAGGER_COMMON
export SWAGGER_V1_ROOT
export SWAGGER_V1_TEMPLATE
export SWAGGER_V1_TEMP_GEN

ifndef PLUGIN_DIR
PLUGIN_DIR := /home/$(USER)/go/plugins
endif
export PLUGIN_DIR

############################
## Declare Module Targets ##
############################
define add_module
	$(eval MODULE := $1)
	$(eval include $(MODULE)/cloud/go/defs.mk)

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
	VET_LIST += $(MODULE)_vet

	MIGRATION_PLUGIN_LIST += $(MODULE)_migrationplugin
endef
_ := $(foreach module, $(MAGMA_MODULES), $(eval $(call add_module,$(module))))


###################
## General Rules ##
###################

precommit: fmt test vet swagger_validate

build: plugin $(BUILD_LIST) migration_plugin swagger
$(BUILD_LIST): %_build:
	make -C $*/cloud/go build

clean: $(CLEAN_LIST)
$(CLEAN_LIST): %_clean:
	make -C $*/cloud/go clean

clean_gen: $(CLEAN_GEN_LIST)
$(CLEAN_GEN_LIST): %_cleangen:
	make -C $*/cloud/go clean_gen

COVER_DIR := /var/tmp
export COVER_DIR
cover: $(COVER_LIST)
$(COVER_LIST): %_cover:
	make -C $*/cloud/go cover

download: $(DOWNLOAD_LIST)
$(DOWNLOAD_LIST): %_download:
	make -C $*/cloud/go download

fmt: $(FMT_LIST)
$(FMT_LIST): %_fmt:
	make -C $*/cloud/go fmt

gen: tools $(GEN_LIST)
$(GEN_LIST): %_gen:
	make -C $*/cloud/go gen

lint: $(LINT_LIST)
$(LINT_LIST): %_lint:
	make -C $*/cloud/go lint

plugin: $(PLUGIN_LIST)
$(PLUGIN_LIST): %_plugin:
	make -C $*/cloud/go plugin

test: tools $(TEST_LIST)
$(TEST_LIST): %_test:
	make -C $*/cloud/go test

tidy: $(TIDY_LIST)
$(TIDY_LIST): %_tidy:
	make -C $*/cloud/go tidy

tools: $(TOOLS_LIST)
$(TOOLS_LIST): %_tools:
	make -C $*/cloud/go tools

vet: $(VET_LIST)
$(VET_LIST): %_vet:
	make -C $*/cloud/go vet

##############################
## Swagger/API docs Targets ##
##############################

swagger: swagger_directories tools $(SWAGGER_LIST)
	cp $(MAGMA_ROOT)/orc8r/cloud/go/models/swagger-common.yml $(SWAGGER_V1_ROOT)/$(SWAGGER_COMMON)
	combine_swagger --inp=$(SWAGGER_V1_TEMP_GEN) --common=$(SWAGGER_V1_ROOT)/$(SWAGGER_COMMON) --out=$(SWAGGER_V1_YML)
	rm -rf $(SWAGGER_V1_TEMP_GEN)

swagger_directories:
	rm -rf $(SWAGGER_V1_TEMP_GEN)
	mkdir $(SWAGGER_V1_TEMP_GEN)

$(SWAGGER_LIST): %_swagger:
	make -C $*/cloud/go copy_swagger_files

swagger_validate: swagger
	swagger validate $(SWAGGER_V1_YML)

##################
## Misc Targets ##
##################

list:  # List all commands
	@echo -e "\nAvailable commands:\n"
	@grep '^[[:lower:]].*:.*#.*' Makefile | awk -F ':|#' '{ printf "\t%-30s%s\n", $$1, $$3 }'
	@echo -e "\n"

migration_plugin: $(MIGRATION_PLUGIN_LIST)
$(MIGRATION_PLUGIN_LIST): %_migrationplugin:
	make -C $*/cloud/go migration_plugin