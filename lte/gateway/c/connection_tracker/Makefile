.PHONY: all build test clean run status

DPI_DIR = $(MAGMA_ROOT)/lte/gateway/c/connection_tracker
BUILD_TYPE ?= Debug

all: build

build: build_connection_tracker

# run_cmake BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, ENV
define run_cmake
mkdir -p $(1)
cd $(1) && $(4) cmake $(2) -DCMAKE_BUILD_TYPE=$(BUILD_TYPE) $(3) -GNinja
ninja -C $(1)
endef

clean:
	rm -rf $(C_BUILD)

build_common:
	$(call run_cmake, $(C_BUILD)/magma_common, $(MAGMA_ROOT)/orc8r/gateway/c/common, )

build_connection_tracker:
	$(call run_cmake, $(C_BUILD)/cont, $(DPI_DIR), )
