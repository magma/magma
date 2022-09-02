default: $(PYTHON_BUILD)/setupinteg_env

include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk
include defs.mk

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
