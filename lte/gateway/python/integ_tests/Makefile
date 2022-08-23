# By default just setup the env.
.PHONY: default
default: $(PYTHON_BUILD)/setupinteg_env

include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk
include defs.mk

# Pull in dependencies from the magma dev build to create virtualenv and
# generate the python protobufs and the common rpc wrappers.
# Also, copy over test config files and update PYTHONPATH to include s1aptester
# lib directory, as well as Swagger-generated API bindings.
$(PYTHON_BUILD)/setupinteg_env:
	make -C .. buildenv
	#@echo "Installing swagger-codegen requirements"
	#$(VIRT_ENV_PIP_INSTALL) -r $(SWAGGER_CODEGEN_OUTPUT)/requirements.txt
	@echo "Copying s1aptester config files"
	cp $(S1AP_TESTER_CFG)/* $(S1AP_TESTER_ROOT)
	@echo "Install MySQL for upstreaming"
	$(VIRT_ENV_PIP_INSTALL) mysqlclient==1.3.13
	@echo "Install flaky for test retries support"
	$(VIRT_ENV_PIP_INSTALL) flaky
	@echo export PYTHONPATH=$(PYTHONPATH):$(S1AP_TESTER_PYTHON_PATH) >> $(PYTHON_BUILD)/bin/activate
	touch $(PYTHON_BUILD)/setupinteg_env

RESULTS_DIR := /var/tmp/test_results
ifdef enable-flaky-retry
	FLAKY_CMD_ARGS := --force-flaky --no-flaky-report --max-runs=3 --min-passes=1
else
	FLAKY_CMD_ARGS :=
endif
# Don't remove the ending comment from function execute_test. Somehow without
# any ending statement the failing testcases are missing status update
define execute_test
	echo "Running test: $(1)"
	timeout --foreground -k 930s 900s sudo -E PATH=$(PATH) PYTHONPATH=$(PYTHONPATH):$(S1AP_TESTER_PYTHON_PATH) $(BIN)/pytest $(FLAKY_CMD_ARGS) --capture=tee-sys --junit-xml=$(RESULTS_DIR)/$(basename $(notdir $(1))).xml -x $(1) || (echo "fail" > $(MAGMA_ROOT)/test_status.txt && exit 1)
	# Testcase $(1) execution completed
endef

.PHONY: precommit
precommit: $(PYTHON_BUILD)/setupinteg_env $(BIN)/pytest
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(PRECOMMIT_TESTS),$(call execute_test,$(test));)
endif

.PHONY: integ_test
integ_test: $(PYTHON_BUILD)/setupinteg_env $(BIN)/pytest
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	echo "pass" > $(MAGMA_ROOT)/test_status.txt
ifndef enable-flaky-retry
	echo "Flaky test retries are disabled"
	$(foreach test,$(EXTENDED_TESTS) $(PRECOMMIT_TESTS),$(call execute_test,$(test));)
else
	echo "Flaky test retries are enabled"
	-$(foreach test,$(EXTENDED_TESTS) $(PRECOMMIT_TESTS),$(call execute_test,$(test));)
	if [ ! -z `grep -s pass $(MAGMA_ROOT)/test_status.txt` ]; then echo "Final integ_test status: Passed"; else echo "Final integ_test status: Failed"; exit 1; fi
endif
endif

.PHONY: fed_integ_test
fed_integ_test: $(PYTHON_BUILD)/setupinteg_env $(BIN)/pytest
	. $(PYTHON_BUILD)/bin/activate
	$(eval export FEDERATED_MODE = True)
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	echo "pass" > $(MAGMA_ROOT)/test_status.txt
ifndef enable-flaky-retry
	echo "Flaky test retries are disabled"
	$(foreach test,$(FEDERATED_TESTS),$(call execute_test,$(test));)
else
	echo "Flaky test retries are enabled"
	-$(foreach test,$(FEDERATED_TESTS),$(call execute_test,$(test));)
	if [ ! -z `grep -s pass $(MAGMA_ROOT)/test_status.txt` ]; then echo "Final fed_integ_test status: Passed"; else echo "Final fed_integ_test status: Failed"; exit 1; fi
endif
endif

.PHONY: nonsanity
nonsanity: $(PYTHON_BUILD)/setupinteg_env $(BIN)/pytest
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(NON_SANITY_TESTS),$(call execute_test,$(test));)
endif

local_integ_test:
	# check if magma services are running
	systemctl is-active --quiet magma@magmad || (echo "Local integ tests should be run on access gw with magma services running"; exit 1)
	. $(PYTHON_BUILD)/bin/activate; sudo $(BIN)/pytest -s $(LOCAL_INTEG_TESTS)
