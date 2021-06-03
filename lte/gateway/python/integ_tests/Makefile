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
	$(VIRT_ENV_PIP_INSTALL) mysqlclient==1.3.12
	@echo export PYTHONPATH=$(PYTHONPATH):$(S1AP_TESTER_PYTHON_PATH) >> $(PYTHON_BUILD)/bin/activate
	touch $(PYTHON_BUILD)/setupinteg_env

# TODO T21489739 - Don't sleep and don't stop after a failure
RESULTS_DIR := /var/tmp/test_results
define execute_test
	echo "Running test: $(1)"
	timeout --foreground -k 930s 900s sudo -E PATH=$(PATH) PYTHONPATH=$(PYTHONPATH):$(S1AP_TESTER_PYTHON_PATH) $(PYTHON_BUILD)/bin/nosetests --with-xunit --xunit-file=$(RESULTS_DIR)/$(basename $(notdir $(1))).xml -x -s $(1) || exit 1
	sleep 1
endef

.PHONY: precommit
precommit: $(PYTHON_BUILD)/setupinteg_env $(BIN)/nosetests
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(PRECOMMIT_TESTS),$(call execute_test,$(test));)
endif

.PHONY: integ_test
integ_test: $(PYTHON_BUILD)/setupinteg_env $(BIN)/nosetests
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(EXTENDED_TESTS) $(PRECOMMIT_TESTS),$(call execute_test,$(test));)
endif

.PHONY: nonsanity
nonsanity: $(PYTHON_BUILD)/setupinteg_env $(BIN)/nosetests
	. $(PYTHON_BUILD)/bin/activate
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(NON_SANITY_TESTS),$(call execute_test,$(test));)
endif

local_integ_test:
	# check if magma services are running
	systemctl is-active --quiet magma@magmad || (echo "Local integ tests should be run on access gw with magma services running"; exit 1)
	. $(PYTHON_BUILD)/bin/activate; sudo $(BIN)/nosetests -s $(LOCAL_INTEG_TESTS)
