# By default do nothing.
.PHONY: default
default:
	# No default set for this Makefile.

include $(MAGMA_ROOT)/orc8r/gateway/python/python.mk
include defs.mk

RESULTS_DIR := /var/tmp/test_results
ifdef enable-flaky-retry
ifneq ($(enable-flaky-retry), false)
	IS_FLAKY_RETRIES_ENABLED := True
	FLAKY_CMD_ARGS := --force-flaky --no-flaky-report --max-runs=3 --min-passes=1
else
	IS_FLAKY_RETRIES_ENABLED := False
	FLAKY_CMD_ARGS :=
endif
else
	IS_FLAKY_RETRIES_ENABLED := False
	FLAKY_CMD_ARGS :=
endif

define apply_final_status
	if [ ! -z `grep -s pass $(MAGMA_ROOT)/test_status.txt` ]; \
	then echo "Final status: Passed"; \
	else echo "Final status: Failed"; exit 1; \
	fi
endef

# Don't remove the ending comment from function execute_test. Somehow without
# any ending statement the failing testcases are missing status update
define execute_test
	echo "Running test: $(1)"
	timeout --foreground -k 930s 900s \
	sudo -E PATH=$(PATH) /tmp/integ_test_env/bin/python3 -m pytest \
	$(FLAKY_CMD_ARGS) \
	--capture=tee-sys \
	--junit-xml=$(RESULTS_DIR)/$(basename $(notdir $(1))).xml -x $(1) \
	|| (echo "fail" > $(MAGMA_ROOT)/test_status.txt && exit 1)
	# Testcase $(1) execution completed
endef

.PHONY: prepare_environment
prepare_environment:
	rm -rf /tmp/integ_test_env
	bazel run --define=on_magma_test=1 //lte/gateway/python/integ_tests/s1aptests:python_env /tmp/integ_test_env
	echo "pass" > $(MAGMA_ROOT)/test_status.txt
	echo "Flaky test retries enabled? $(IS_FLAKY_RETRIES_ENABLED)"

.PHONY: prepare_federation
prepare_federation:
	$(eval export FEDERATED_MODE = True)

.PHONY: selected_tests
selected_tests: prepare_environment
	-$(foreach test,$(TESTS),$(call execute_test,$(test));)

.PHONY: precommit
precommit: prepare_environment
	-$(foreach test,$(PRECOMMIT_TESTS),$(call execute_test,$(test));)

.PHONY: extended_tests
extended_tests: prepare_environment
	-$(foreach test,$(EXTENDED_TESTS),$(call execute_test,$(test));)

.PHONY: extended_tests_long
extended_tests_long: prepare_environment
	-$(foreach test,$(EXTENDED_TESTS_LONG),$(call execute_test,$(test));)

.PHONY: integ_test
integ_test: precommit extended_tests extended_tests_long

.PHONY: integ_test_containerized
integ_test_containerized: precommit extended_tests extended_tests_long

.PHONY: federated_integ_test
federated_integ_test: prepare_environment prepare_federation
	-$(foreach test,$(FEDERATED_TESTS),$(call execute_test,$(test));)

.PHONY: nonsanity
nonsanity: prepare_environment
	-$(foreach test,$(NON_SANITY_TESTS),$(call execute_test,$(test));)

.PHONY: evaluate_result
evaluate_result:
	$(call apply_final_status)
