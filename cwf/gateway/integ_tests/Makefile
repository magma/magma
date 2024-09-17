.PHONY: integ_test

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

# In Go, you can specify tests to run by prefix
# Integration tests to run (Tests with a prefix name Test)
MANDATORY_TESTS=Test
AUTH=Authenticate
GX=Gx
GY=Gy
MULTISP=MultiSessionProxy
HSSLESS=Hssless

define execute_test
 	echo "Running test: $(1)"
 	# tag  main_integration prevents running multi_session_proxy tests.
 	# multi_session_proxy_must be run excusively from multi_session_proxy target
	go test -tags main_integration -run $(1)
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

multi_session_proxy:
	echo "Running Multi Session Proxy Tests..."
	go test -run $(MULTISP)

hssless:
	echo "Running HSS Less Tests..."
	$(foreach test,$(HSSLESS),$(call execute_test,$(test));)

restart:
	echo "Running restart tests..."
	go test -tags=restart -timeout 20m

precommit:
	go fmt ./...
	go test --tags=all -c ./...
