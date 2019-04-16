.PHONY: all build run

PIP_INSTALL := $(PYTHON_BUILD)/bin/pip3 install -q -U --cache-dir $(PIP_CACHE_HOME)

build:
	$(MAKE) -C ../../orc8r/gateway/python buildenv
	$(PIP_INSTALL) -e python

run: build
	sudo service magma@magmad start
