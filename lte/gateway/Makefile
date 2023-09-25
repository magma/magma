.PHONY: all build clean run test

define DEPRECATION_WARNING

Warning: As part of the Make to Bazel switch-over of the AGW this Makefile
will be deprecated soon, see https://github.com/magma/magma/issues/14971.
Do not add any new Make targets for the AGW!
All new targets should be integrated with the Bazel build system,
see https://magma.github.io/magma/docs/next/bazel/agw_with_bazel.

endef
export DEPRECATION_WARNING

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

TEST_FLAG = -DBUILD_TESTS=1
OAI_NOTEST_FLAGS = -DMME_UNIT_TEST=False

all: build

build: build_python build_common build_oai build_sctpd build_session_manager build_connection_tracker build_envoy_controller build_li_agent ## Build all

clean: clean_python clean_envoy_controller ## Clean all builds
	@echo "$$DEPRECATION_WARNING"
	rm -rf $(C_BUILD)

clean_python: ## Clean Python-only builds
	@echo "$$DEPRECATION_WARNING"
	make -C $(MAGMA_ROOT)/lte/gateway/python clean

clean_envoy_controller: ## Clean envoy controller build
	@echo "$$DEPRECATION_WARNING"
	rm -rf  $(GO_BUILD)/envoy_controller

run: build ## Build and run all services
	@echo "$$DEPRECATION_WARNING"
	sudo service magma@* stop
	sudo service magma@magmad start

# run_cmake BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, ENV
define run_cmake
mkdir -p $(1)
cd $(1) && $(4) cmake $(2) -DCMAKE_BUILD_TYPE=$(BUILD_TYPE) $(3) -GNinja
# Turn on generation of compile_commands.json (See GH Issue #5488)
cd $(1) && awk '{if (/^CMAKE_EXPORT_COMPILE_COMMANDS/) gsub(/OFF/, "ON"); print}' CMakeCache.txt > tmp && mv tmp CMakeCache.txt
ninja -C $(1)
endef

# run_ctest BUILD_DIRECTORY, TEST_BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, LIST OF TESTS
define run_ctest
$(call run_cmake, $(1), $(3), $(4) $(TEST_FLAG))
cd $(2) && ctest --output-on-failure -R $(5)
endef

build_python: ## Build Python environment
	@echo "$$DEPRECATION_WARNING"
	sudo service magma@* stop
	make -C $(MAGMA_ROOT)/lte/gateway/python buildenv

build_common: ## Build shared libraries
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/magma_common, $(MAGMA_ROOT)/orc8r/gateway/c/common, $(COMMON_FLAGS))

build_oai: ## Build OAI
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(COMMON_FLAGS) $(OAI_NOTEST_FLAGS))

build_sctpd:
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/sctpd, $(GATEWAY_C_DIR)/sctpd, )

build_session_manager: build_common ## Build session manager
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/session_manager, $(GATEWAY_C_DIR)/session_manager, )

build_li_agent: ## Build li agent
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/li_agent, $(GATEWAY_C_DIR)/li_agent, )

build_connection_tracker:
	@echo "$$DEPRECATION_WARNING"
	$(call run_cmake, $(C_BUILD)/connection_tracker, $(GATEWAY_C_DIR)/connection_tracker, )

build_envoy_controller: ## Build envoy controller
	@echo "$$DEPRECATION_WARNING"
	cd $(MAGMA_ROOT)/feg/gateway && $(MAKE) install_envoy_controller

# Catch all for c service tests
# This works with test_dpi and test_session_manager
test_%: build_common
	@echo "$$DEPRECATION_WARNING"
	$(call run_ctest, $(C_BUILD)/$*, $(C_BUILD)/$*, $(GATEWAY_C_DIR)/$*, )
