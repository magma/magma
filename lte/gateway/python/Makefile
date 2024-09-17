include $(MAGMA_ROOT)/lte/gateway/python/defs.mk
# Contains common targets for gateway python development
include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk

define DEPRECATION_WARNING

Warning: As part of the Make to Bazel switch-over of the AGW this Makefile
will be deprecated soon, see https://github.com/magma/magma/issues/14971.
Do not add any new Make targets for the AGW!
All new targets should be integrated with the Bazel build system,
see https://magma.github.io/magma/docs/next/bazel/agw_with_bazel.

endef
export DEPRECATION_WARNING

# Set up some pattern rules
define setup_rules
	$(eval PSRC := $1)
	CLEAN_LIST += $(PSRC)_clean
	BUILD_LIST += $(PSRC)_build
	TEST_LIST += $(PSRC)_test
endef
_ := $(foreach python_src, $(PYTHON_SRCS), $(eval $(call setup_rules,$(python_src))))

snowflake:
	@echo "$$DEPRECATION_WARNING"
	sudo env PATH=${PYTHON_BUILD}/bin:${PATH} snowflake --make-snowflake

buildenv: setupenv protos swagger $(BUILD_LIST) py_patches snowflake
$(BUILD_LIST): %_build:
	@echo "$$DEPRECATION_WARNING"
	make -C $* install_egg

test_all: buildenv $(BIN)/pytest $(BIN)/pytest-cov $(TEST_LIST)
	@echo "$$DEPRECATION_WARNING"
$(TEST_LIST): %_test:
	make -C $* .test

define run_unit_tests
    . $(PYTHON_BUILD)/bin/activate; cd $(2); $(PYTHON_BUILD)/bin/pytest -s $(1) || exit 1
endef
define run_sudo_unit_tests
    . $(PYTHON_BUILD)/bin/activate; cd $(2); sudo $(PYTHON_BUILD)/bin/pytest -s $(1) || exit 1
endef

# unit_tests (UT_PATH=unit_test_path|MAGMA_SERVICE=service_name)[DONT_BUILD_ENV=1]
unit_tests: $(TESTS)
	@echo "$$DEPRECATION_WARNING"
ifndef MAGMA_SERVICE
ifndef UT_PATH
	@echo "usage: make unit_tests (UT_PATH=unit_test_path|MAGMA_SERVICE=service_name)[DONT_BUILD_ENV=1]"
	@exit 1
endif
endif
ifdef MAGMA_SERVICE
	$(eval SELECTED_TESTS ?= $(shell grep -oP '[^\s]+(?=$(MAGMA_SERVICE))[^\s]+' $(MAGMA_ROOT)/lte/gateway/python/defs.mk))
	$(eval TEST_PATH_ORC8R ?= $(shell grep -oP '[^\s]+(?=$(MAGMA_SERVICE))[^\s]+' $(MAGMA_ROOT)/orc8r/gateway/python/defs.mk))
else ifdef UT_PATH
	$(eval SELECTED_TESTS ?= $(subst $(MAGMA_ROOT)/lte/gateway/python/, , $(UT_PATH)))
endif
	@if [ ! -d "$(MAGMA_ROOT)/lte/gateway/python/$(TEST_PATH)" ]; then if [ ! -d "$(MAGMA_ROOT)/orc8r/gateway/python/$(TEST_PATH_ORC8R)" ]; then echo "no tests found" && exit 1; fi; fi
ifndef DONT_BUILD_ENV
	@$(MAKE) buildenv $(BIN)/pytest
endif
ifdef MAGMA_SERVICE
	sudo service magma@$(MAGMA_SERVICE) stop
endif
	$(eval NON_SUDO_TESTS ?= $(patsubst %, \%%,$(TESTS)))
	$(eval M_SUDO_TESTS ?= $(patsubst %, \%%,$(SUDO_TESTS)))

	$(if $(strip $(NON_SUDO_TESTS)),$(eval SEL_TEST_PATH_LTE ?= $(filter $(NON_SUDO_TESTS), $(SELECTED_TESTS))))
	$(if $(strip $(M_SUDO_TESTS)),$(eval SEL_SUDO_TEST_PATH_LTE ?= $(filter $(M_SUDO_TESTS), $(SELECTED_TESTS))))

	$(if $(strip $(SEL_TEST_PATH_LTE)),$(call run_unit_tests,$(SEL_TEST_PATH_LTE), $(MAGMA_ROOT)/lte/gateway/python/))
	$(if $(strip $(SEL_SUDO_TEST_PATH_LTE)),$(call run_sudo_unit_tests,$(SEL_SUDO_TEST_PATH_LTE), $(MAGMA_ROOT)/lte/gateway/python/))
	$(if $(strip $(TEST_PATH_ORC8R)), $(call run_unit_tests, $(TEST_PATH_ORC8R), $(MAGMA_ROOT)/orc8r/gateway/python/))

$(BIN)/pytest: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) pytest==7.1.2

$(BIN)/pylint: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) pylint==2.14.0

# Disable W0511: todo warnings
# Disable R0903: Too few public methods

check:
	./precommit.py --format --diff

clean: $(CLEAN_LIST)
	@echo "$$DEPRECATION_WARNING"
	sudo rm -rf $(PYTHON_BUILD)/
	sudo find . -name '*.pyc' -o -name '__pycache__' -prune -exec rm -rf {} \;
$(CLEAN_LIST): %_clean:
	make -C $* remove_egg
