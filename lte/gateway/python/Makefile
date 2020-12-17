include $(MAGMA_ROOT)/lte/gateway/python/defs.mk
# Contains common targets for gateway python development
include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk

# Set up some pattern rules
define setup_rules
	$(eval PSRC := $1)
	CLEAN_LIST += $(PSRC)_clean
	BUILD_LIST += $(PSRC)_build
	TEST_LIST += $(PSRC)_test
endef
_ := $(foreach python_src, $(PYTHON_SRCS), $(eval $(call setup_rules,$(python_src))))

buildenv: setupenv protos swagger $(BUILD_LIST) py_patches
$(BUILD_LIST): %_build:
	make -C $* install_egg

test_all: buildenv $(BIN)/nosetests $(BIN)/coverage $(TEST_LIST)
$(TEST_LIST): %_test:
	make -C $* .test

coverage: $(BIN)/coverage
	$(BIN)/coverage report

$(BIN)/nosetests: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) -I nose==1.3.7

$(BIN)/coverage: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) "coverage==4.5.1"

$(BIN)/pylint: install_virtualenv
    # https://github.com/PyCQA/pylint/issues/380
    # compile issue with py3.4 so install with no-compile
	$(VIRT_ENV_PIP_INSTALL) pylint==1.6.4 --no-compile

$(BIN)/pep8: install_virtualenv
	# pylint doesn't cover all the pep8 style guidelines. Specifically,
	# E203, E301, E303, W203, W291, W292
	$(VIRT_ENV_PIP_INSTALL) pep8==1.7.0

# Disable C0330: bad-continuation errors with pylint1.6.4 as it conflicts with pep8
# 	https://github.com/pylint-bot/pylint-unofficial/issues/289
# Disable W0511: todo warnings
# Disable R0903: Too few public methods
CHECK_CMD_PYLINT := find . -name '*.py' -exec $(BIN)/pylint --disable=C0330,R0903,W0511 {} +;
CHECK_CMD_PEP8 := find . -name '*.py' -exec $(BIN)/pep8 {} +;
check: buildenv $(BIN)/pylint $(BIN)/pep8
	$(CHECK_CMD_PEP8)
	$(CHECK_CMD_PYLINT)

clean: $(CLEAN_LIST)
	sudo rm -rf $(PYTHON_BUILD)/ .coverage
	sudo find . -name '*.pyc' -o -name '__pycache__' -prune -exec rm -rf {} \;
$(CLEAN_LIST): %_clean:
	make -C $* remove_egg
