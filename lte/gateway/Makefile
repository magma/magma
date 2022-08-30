.PHONY: all build clean log logs run status test

GATEWAY_C_DIR = $(MAGMA_ROOT)/lte/gateway/c
GRPC_CPP_PLUGIN_PATH ?= `which grpc_cpp_plugin`
BUILD_TYPE ?= Debug

# FEATURES: What kind of flavours do you want your MME or AGW have in it
# MME is MME as described in 3GPP specs, it has at least S1AP, S11, S6a
# interfaces.
# AGW means Acces GateWay, is the result of the aggregation of MME, SGW and PGW.
# First in FEATURES, select what to you want to build : mme or agw with OpenFlow
# (OVS): FEATURE=mme_oai or agw_of
# Then you can have other features that can be built for mme or agw :
# s6a with fd (freeDiameter)

# Default is agw with OpenFlow, gRPC over S6a , (no freeDiameter over s6a).
FEATURES ?= agw_of
# AVAILABLE_FEATURE_LIST : every feature not in this list will trigger an error.
AVAILABLE_FEATURE_LIST = agw_of mme_oai
REQUESTED_FEATURE_LIST = $(sort $(FEATURES))

# First, check that nothing outside of AVAILABLE_FEATURE_LIST is requested
ifneq ($(words $(strip $(filter-out $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))), 0)
  $(error Non allowed flags: "$(filter-out $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST))")
endif

# Then check that only one of the available options are selected
ifneq ($(words $(strip $(filter $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))), 1)
	$(error Exclusive flags: "$(filter $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST))")
endif

MAIN_FEATURE = $(strip $(filter $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))
$(info MAIN_FEATURE $(MAIN_FEATURE))

ifeq ($(MAIN_FEATURE),mme_oai)
# Set DS6A_OVER_GRPC=False if using OAI-HSS
OAI_FLAGS = -DS6A_OVER_GRPC=False -DEMBEDDED_SGW=False
else ifeq ($(MAIN_FEATURE),agw_of)
OAI_FLAGS = -DS6A_OVER_GRPC=True -DEMBEDDED_SGW=True
else
# Error if nothing is specified
$(error No feature is specified! This should not happen)
endif

# debian stretch build uses older cc not recognizing options needed on ubuntu focal

OS_VERSION_NAME := $(shell (grep VERSION_CODENAME /etc/os-release || true) | sed 's/.*=//g')
COMMON_FLAGS = -DCMAKE_C_FLAGS="-Wall $(CPPFLAGS)" -DCMAKE_CXX_FLAGS="-Wall $(CPPFLAGS)"

$(info OAI_FLAGS $(OAI_FLAGS))

FUZZ_FLAGS = $(OAI_FLAGS) -DFUZZ=True
TEST_FLAG = -DBUILD_TESTS=1
OAI_TEST_FLAGS = -DMME_UNIT_TEST=True
OAI_NOTEST_FLAGS = -DMME_UNIT_TEST=False
OAI_TESTS ?= ".*"
OAI_BENCHMARK_FLAGS = -DMME_BENCHMARK=True
OAI_NOBENCHMARK_FLAGS = -DMME_BENCHMARK=False

all: build

build: build_python build_common build_oai build_sctpd build_session_manager build_connection_tracker build_envoy_controller build_li_agent ## Build all

smf_build: build_session_manager  ## Build only sessionD component make smf_build

test: test_python test_common test_oai test_sctpd test_session_manager ## Run all tests

clean: clean_python clean_envoy_controller ## Clean all builds
	rm -rf $(C_BUILD)

clean_python: ## Clean Python-only builds
	make -C $(MAGMA_ROOT)/lte/gateway/python clean

clean_envoy_controller: ## Clean envoy controller build
	rm -rf  $(GO_BUILD)/envoy_controller

start: ## Start all services
	sudo service magma@magmad start

stop: ## Stop all services
	sudo service magma@* stop

restart: stop start ## Restart all services

run: build restart ## Build and run all services

status: ## Status of all services
	sudo service magma@* status

log: ## Follow logs for magmad service
	sudo journalctl -fu magma@magmad | egrep 'error|$$' -i --color

logs: ## Follow logs for all services
	sudo journalctl -fu magma@* | egrep 'error|$$' -i --color

# run_cmake BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, ENV
define run_cmake
mkdir -p $(1)
cd $(1) && $(4) cmake $(2) -DCMAKE_BUILD_TYPE=$(BUILD_TYPE) $(3) -GNinja
# Turn on generation of compile_commands.json (See GH Issue #5488)
cd $(1) && awk '{if (/^CMAKE_EXPORT_COMPILE_COMMANDS/) gsub(/OFF/, "ON"); print}' CMakeCache.txt > tmp && mv tmp CMakeCache.txt
ninja -C $(1)
endef

# run_scanbuild BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS
define run_scanbuild
$(eval REPORT_DIR = "$(1)/reports")
mkdir -p $(1)
mkdir -p $(REPORT_DIR)
cd $(1) && scan-build cmake $(2) -DCMAKE_BUILD_TYPE=Debug $(3) -GNinja
scan-build -o $(REPORT_DIR) ninja -C $(1)
cp -r $(REPORT_DIR) $(MAGMA_ROOT)
@echo "Reports in magma/reports/.../index.html"
endef

# run_ctest BUILD_DIRECTORY, TEST_BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, LIST OF TESTS
define run_ctest
$(call run_cmake, $(1), $(3), $(4) $(TEST_FLAG))
cd $(2) && ctest --output-on-failure -R $(5)
endef

build_python: stop ## Build Python environment
	make -C $(MAGMA_ROOT)/lte/gateway/python buildenv

build_common: ## Build shared libraries
	$(call run_cmake, $(C_BUILD)/magma_common, $(MAGMA_ROOT)/orc8r/gateway/c/common, $(COMMON_FLAGS))

benchmark_pb: ## Benchmark ProtoBuf, same build_oai but with -DMME_BENCHMARK=True
	$(call run_cmake, $(C_BUILD)/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(COMMON_FLAGS) $(OAI_NOTEST_FLAGS) $(OAI_BENCHMARK_FLAGS))

build_oai: ## Build OAI
	$(call run_cmake, $(C_BUILD)/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(COMMON_FLAGS) $(OAI_NOTEST_FLAGS) $(OAI_NOBENCHMARK_FLAGS))

format_all:
	find $(MAGMA_ROOT)/orc8r/gateway/c/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h" -o -iname "*.hpp" \) -exec \
	/usr/bin/clang-format-11 -i {} \;
	find $(MAGMA_ROOT)/lte/gateway/c/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h"  -o -iname "*.hpp" \) -exec \
	/usr/bin/clang-format-11 -i {} \;
	find $(MAGMA_ROOT)/lte/gateway/python/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h"  -o -iname "*.hpp" \) -exec \
	/usr/bin/clang-format-11 -i {} \;

build_sctpd:
	$(call run_cmake, $(C_BUILD)/sctpd, $(GATEWAY_C_DIR)/sctpd, )

build_session_manager: build_common ## Build session manager
	$(call run_cmake, $(C_BUILD)/session_manager, $(GATEWAY_C_DIR)/session_manager, )

build_li_agent: ## Build li agent
	$(call run_cmake, $(C_BUILD)/li_agent, $(GATEWAY_C_DIR)/li_agent, )

build_connection_tracker:
	$(call run_cmake, $(C_BUILD)/connection_tracker, $(GATEWAY_C_DIR)/connection_tracker, )

build_envoy_controller: ## Build envoy controller
	cd $(MAGMA_ROOT)/feg/gateway && $(MAKE) install_envoy_controller

# Catch all for c services that don't have custom flags
# This works with build_dpi
build_%:
	$(call run_cmake, $(C_BUILD)/$*, $(MAGMA_ROOT)/c/$*, $(COMMON_FLAGS))

scan_oai: ## Scan OAI
	$(call run_scanbuild, $(C_BUILD)/scan/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS))

test_python: stop ## Run all Python-specific tests
	make -C $(MAGMA_ROOT)/lte/gateway/python test_all

test_sudo_python: stop ## Run Python tests that require sudo (datapath, etc.)
	make -C $(MAGMA_ROOT)/lte/gateway/python test_all SKIP_NON_SUDO_TESTS=1

test_python_service: ## Run all Python-specific tests for a given service
ifdef UT_PATH
	$(eval ut_path?=$(shell realpath $(UT_PATH)))
endif
	make -C $(MAGMA_ROOT)/lte/gateway/python unit_tests MAGMA_SERVICE=$(MAGMA_SERVICE) UT_PATH=$(ut_path) DONT_BUILD_ENV=$(DONT_BUILD_ENV)

test_oai: ## Run all OAI-specific tests
	$(call run_ctest, $(C_BUILD)/core, $(C_BUILD)/core/oai, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(OAI_TEST_FLAGS) $(OAI_NOBENCHMARK_FLAGS), $(OAI_TESTS))

test_oai_runtime: export REPORT_FOLDER=${MAGMA_ROOT}/report/
test_oai_runtime: export UNITTEST_REPORT_FOLDER=${REPORT_FOLDER}/unittest_report/
test_oai_runtime: export MERGED_REPORT_FOLDER=${REPORT_FOLDER}/merged_report/
test_oai_runtime: export GTEST_OUTPUT=xml:${UNITTEST_REPORT_FOLDER}
test_oai_runtime: ## Run all OAI-specific tests with report about the running time
	mkdir -p ${UNITTEST_REPORT_FOLDER}
	mkdir -p ${MERGED_REPORT_FOLDER}
	-$(call run_ctest, $(C_BUILD)/core, $(C_BUILD)/core/oai, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(OAI_TEST_FLAGS) $(OAI_NOBENCHMARK_FLAGS), $(OAI_TESTS))
	python3 python/scripts/runtime_report.py -i .+\\.xml$$ -w ${UNITTEST_REPORT_FOLDER}
	rm ${UNITTEST_REPORT_FOLDER}/*.xml

test_sctpd: ## Run all tests for sctpd
	$(call run_ctest, $(C_BUILD)/sctp, $(C_BUILD)/sctp/src, $(GATEWAY_C_DIR)/sctpd, )

test_common: ## Run all tests in magma_common
	$(call run_cmake, $(C_BUILD)/magma_common, $(MAGMA_ROOT)/orc8r/gateway/c/common, $(TEST_FLAG))
	# Run the common lib tests that exist
	cd $(C_BUILD)/magma_common/config && ctest --output-on-failure
	cd $(C_BUILD)/magma_common/service303 && ctest --output-on-failure
	cd $(C_BUILD)/magma_common/service_registry && ctest --output-on-failure

test_li_agent:
	$(call run_ctest, $(C_BUILD)/li_agent, $(C_BUILD)/li_agent/src, $(GATEWAY_C_DIR)/li_agent, )


# Catch all for c service tests
# This works with test_dpi and test_session_manager
test_%: build_common
	$(call run_ctest, $(C_BUILD)/$*, $(C_BUILD)/$*, $(GATEWAY_C_DIR)/$*, )

coverage_oai: test_oai
	lcov --capture --directory $(C_BUILD) --output-file /tmp/coverage_oai.info.raw
	lcov -r /tmp/coverage_oai.info.raw "/*/test/*" "/usr/*" "/build/*protos*" -o /tmp/coverage_oai.info
	rm -f `find $(C_BUILD) -name *.gcda` # Clean up any prior coverage data

# format and test c/session_manager
precommit_sm: format_all test_session_manager

# format and test c/oai
precommit_oai: format_all test_oai

## Generate complete code structural information prior to any test execution
base_coverage: build_oai
	lcov --initial --directory $(C_BUILD) -c --output-file /tmp/coverage_initialize.info.raw
	lcov -r /tmp/coverage_initialize.info.raw "/*/test/*" "/usr/*" "/build/*protos*" -o /tmp/coverage_initialize.info
	rm -f `find $(C_BUILD) -name *.gcda` # Clean up any prior coverage data

# Combine results of sub-coverages
coverage: base_coverage coverage_oai
	lcov -a /tmp/coverage_initialize.info -a /tmp/coverage_oai.info -o $(C_BUILD)/coverage.info
