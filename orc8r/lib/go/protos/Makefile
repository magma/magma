# Targets for protos submodule
clean::
	go clean ./...

fmt::
	go fmt ./...

gen::
	go generate ./...

test::
	$(eval NAME ?= $(shell pwd | tr / _))
	gotestsum --junitfile $(TEST_RESULTS_DIR)/$(NAME).xml

vet::
	go vet -composites=false ./...
