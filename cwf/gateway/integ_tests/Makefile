.PHONY: integ_test

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

# In Go, you can specify tests to run by prefix
# Integration tests to run (Tests with a prefix name Test)
MANDATORY_TESTS=Test
AUTH=Authenticate
GX=Authenticate Gx
GY=Authenticate Gy

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

authenticate:
	echo "Running Auth Tests..."
	$(foreach test,$(AUTH),$(call execute_test,$(test));)

gx:
	echo "Running Auth+Gx Tests..."
	$(foreach test,$(GX),$(call execute_test,$(test));)

gy:
	echo "Running Auth+Gy Tests..."
	$(foreach test,$(GY),$(call execute_test,$(test));)