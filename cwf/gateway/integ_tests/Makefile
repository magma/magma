.PHONY: integ_test

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

# Integration tests to run (Tests with a prefix name TestAuthenticate)  
MANDATORY_TESTS=Test

define execute_test
 	echo "Running test: $(1)"
	go test -run $(1)
endef

integ_test:
ifdef TESTS
	$(call execute_test,$(TESTS))
else
	$(foreach test,$(MANDATORY_TESTS),$(call execute_test,$(test));)
endif
