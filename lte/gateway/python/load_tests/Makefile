include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk
include defs.mk

load_test_file = $(firstword $(subst :, ,$1))
method = $(or $(word 2,$(subst :, ,$1)))

define execute_test
	echo "Running load test: $(1)"
	timeout --foreground -k 930s 900s sudo -E PATH=$(PATH) $(PYTHON_BUILD)/bin/python3 $(call load_test_file,$(1)) $(call method,$(1)) || exit 1
	sleep 1
endef

.PHONY: load_test
load_test:
	. $(PYTHON_BUILD)/bin/activate
	$(foreach test,$(LOAD_TESTS) ,$(call execute_test,$(test));)
