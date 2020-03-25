# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
###############################################################################
# This file contains common Make targets related to setting up a Python
# environment, running tests, and cleaning up. See lte/gateway/python/Makefile
# for an example of how to use this file.
###############################################################################

# virtualenv bin and build dirs
PYTHON_VERSION=3.5
BIN := $(PYTHON_BUILD)/bin
SRC := $(MAGMA_ROOT)
SITE_PACKAGES_DIR := $(PYTHON_BUILD)/lib/python$(PYTHON_VERSION)/site-packages

# Command to pip install into the virtualenv
VIRT_ENV_PIP_INSTALL := $(BIN)/pip3 install -q -U --cache-dir $(PIP_CACHE_HOME)

install_virtualenv:
	@echo "Initializing virtualenv with python version $(PYTHON_VERSION)"
	virtualenv --system-site-packages --python=/usr/bin/python$(PYTHON_VERSION) $(PYTHON_BUILD)
	. $(PYTHON_BUILD)/bin/activate;
	$(VIRT_ENV_PIP_INSTALL) "pip>=19.1.1"

setupenv: $(PYTHON_BUILD)/sysdeps $(SITE_PACKAGES_DIR)/setuptools

# Sytem packages needed for build
SYS_DEPENDENCIES := python3-dev
$(PYTHON_BUILD)/sysdeps: $(PYTHON_BUILD)
	sudo apt-get -y install $(SYS_DEPENDENCIES)
	touch $(PYTHON_BUILD)/sysdeps

$(PYTHON_BUILD):
	mkdir -p $(PYTHON_BUILD)

$(SITE_PACKAGES_DIR)/setuptools: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) "setuptools>=41.0.1"

swagger:: swagger_prereqs $(SWAGGER_LIST)
swagger_prereqs:
	test -f /usr/bin/java # Java exists
	test -n "$(SWAGGER_CODEGEN_JAR)" # SWAGGER_CODEGEN_JAR set
	test -f $(SWAGGER_CODEGEN_JAR) # swagger-codegen exists
	@mkdir -p $(PYTHON_BUILD)/gen
$(SWAGGER_LIST): %_swagger_specs:
	@echo "Generating python code for $* swagger*.yml files"
	@# Clean directory for easy moving of files
	@rm -rf $(PYTHON_BUILD)/gen/$*/swagger
	@mkdir -p $(PYTHON_BUILD)/gen/$*/swagger
	@touch $(PYTHON_BUILD)/gen/$*/swagger/__init__.py
	@# Initialize a subdirectory to store swagger specs
	@mkdir -p $(PYTHON_BUILD)/gen/$*/swagger/specs
	@touch $(PYTHON_BUILD)/gen/$*/swagger/specs/__init__.py
	@# Copy swagger specs over to the build directory,
	@# so that eventd can access them at runtime
	cp $(SRC)/$*/swagger/*.yml $(PYTHON_BUILD)/gen/$*/swagger/specs
	@# Generate the files
	ls $(PYTHON_BUILD)/gen/$*/swagger/specs \
		| grep -e ".*\.yml" \
		| xargs -t -I% /usr/bin/java -jar "$(SWAGGER_CODEGEN_JAR)" generate \
			-i $(PYTHON_BUILD)/gen/$*/swagger/specs/% \
			-o $(PYTHON_BUILD)/gen/$*/swagger \
			-l python \
			-Dmodels
	@# Flatten and clean up directory
	@mv $(PYTHON_BUILD)/gen/$*/swagger/swagger_client/* $(PYTHON_BUILD)/gen/$*/swagger/
	@rmdir $(PYTHON_BUILD)/gen/$*/swagger/swagger_client
	@rm -r $(PYTHON_BUILD)/gen/$*/swagger/test

protos:: $(BIN)/grpcio-tools $(PROTO_LIST) prometheus_proto
	@find $(PYTHON_BUILD)/gen -type d | tail -n +2 | sed '/__pycache__/d' | xargs -I % touch "%/__init__.py"
$(PROTO_LIST): %_protos:
	@echo "Generating python code for $* .proto files"
	@mkdir -p $(PYTHON_BUILD)/gen
	@echo "$(PYTHON_BUILD)/gen" > $(SITE_PACKAGES_DIR)/magma_gen.pth
	$(BIN)/python $(SRC)/protos/gen_protos.py $(SRC)/$*/protos/ $(MAGMA_ROOT),$(MAGMA_ROOT)/orc8r/protos/prometheus $(SRC) $(PYTHON_BUILD)/gen/

prometheus_proto:
	$(BIN)/python $(SRC)/protos/gen_prometheus_proto.py $(MAGMA_ROOT) $(PYTHON_BUILD)/gen

# If you update the version here, you probably also want to update it in setup.py
$(BIN)/grpcio-tools: install_virtualenv
	$(VIRT_ENV_PIP_INSTALL) "grpcio-tools==1.16.1"

.test: .tests .sudo_tests

.tests:
ifdef TESTS
	. $(PYTHON_BUILD)/bin/activate; $(BIN)/nosetests --with-coverage --cover-erase --cover-branches --cover-package=magma -s $(TESTS)
endif

.sudo_tests:
ifdef SUDO_TESTS
ifndef SKIP_SUDO_TESTS
	. $(PYTHON_BUILD)/bin/activate; sudo $(BIN)/nosetests --with-coverage --cover-branches --cover-package=magma -s $(SUDO_TESTS)
endif
endif

install_egg: install_virtualenv setup.py
	$(eval NAME ?= $(shell $(BIN)/python setup.py --name))
	@echo "Installing egg link for $(NAME)"
	$(VIRT_ENV_PIP_INSTALL) -e .[dev]

remove_egg:
	rm -rf *.egg-info
