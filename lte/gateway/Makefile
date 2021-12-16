.PHONY: all build clean help log logs run status test

GATEWAY_C_DIR = $(MAGMA_ROOT)/lte/gateway/c
GRPC_CPP_PLUGIN_PATH ?= `which grpc_cpp_plugin`
BUILD_TYPE ?= Debug
ENABLE_ASAN ?= 0

# FEATURES: What kind of flavours do you want your MME or AGW have in it
# MME is MME as described in 3GPP specs, it has at least S1AP, S11, S6a
# interfaces.
# AGW means Acces GateWay, is the result of the aggregation of MME, SGW and PGW.
# First in FEATURES, select what to you want to build : mme or agw with OpenFlow
# (OVS): FEATURE=mme or agw_of
# Then you can have other features that can be built for mme or agw :
# s6a with fd (freeDiameter)

# Default is agw with OpenFlow, gRPC over S6a , (no freeDiameter over s6a).
FEATURES ?= agw_of
# EXCLUSIVE_FEATURE_LIST : list of primary features that cannot be requested
# together.
EXCLUSIVE_FEATURE_LIST = mme agw_of mme_oai
# AVAILABLE_FEATURE_LIST : every feature not in this list will trigger an error.
AVAILABLE_FEATURE_LIST = s6a_fd mme agw_of mme_oai
REQUESTED_FEATURE_LIST = $(sort $(FEATURES))

ifeq ($(BUILD_TYPE),Debug)
	BAZEL_FLAGS = --config=development
	ifeq ($(ENABLE_ASAN),1)
		BAZEL_FLAGS := $(BAZEL_FLAGS) --config=asan
	endif
else 
	## RelWithDebInfo enable LSAN and add debug information
	BAZEL_FLAGS = --config=production
endif
$(info BAZEL_FLAGS $(BAZEL_FLAGS))

ifneq ($(words $(strip $(filter-out $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))), 0)
  $(error Non allowed flags: "$(filter-out $(AVAILABLE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST))")
endif

ifneq ($(words $(strip $(filter $(EXCLUSIVE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))), 0)
  ifneq ($(words $(strip $(filter $(EXCLUSIVE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))), 1)
    $(error Exclusive flags: "$(filter $(EXCLUSIVE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST))")
  endif
endif

MAIN_FEATURE = $(strip $(filter $(EXCLUSIVE_FEATURE_LIST),$(REQUESTED_FEATURE_LIST)))

ifneq ($(words $(strip $(filter $(REQUESTED_FEATURE_LIST),s6a_fd))), 0)
S6A_FLAGS = -DS6A_OVER_GRPC=False
else
#default
S6A_FLAGS = -DS6A_OVER_GRPC=True
endif

ifeq ($(MAIN_FEATURE),mme)
# Here force S6A to use freeDiameter
OAI_FLAGS = -DS6A_OVER_GRPC=True
else ifeq ($(MAIN_FEATURE),mme_oai)
# Set DS6A_OVER_GRPC=False if using OAI-HSS
OAI_FLAGS = -DS6A_OVER_GRPC=False
else ifeq ($(MAIN_FEATURE),agw_of)
OAI_FLAGS = $(S6A_FLAGS) -DEMBEDDED_SGW=True

else
# default
OAI_FLAGS = $(S6A_FLAGS) -DEMBEDDED_SGW=True

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

# Ref: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show documented commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

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

define copy_bazel_c_build ## Copy Bazel build output to C_BUILD/
# 1 - source directory, 2 - binary name
mkdir -p $(C_BUILD)/$(1)
sudo cp -f $(MAGMA_ROOT)/bazel-bin/lte/gateway/c/$(1)/$(2) $(C_BUILD)/$(1)/$(2)
endef

build_c:  ## Build C/C++ targets with Bazel
	bazel build $(BAZEL_FLAGS) //lte/gateway/c/session_manager:sessiond //lte/gateway/c/sctpd/src:sctpd //lte/gateway/c/connection_tracker/src:connectiond //lte/gateway/c/li_agent/src:liagentd
	$(call copy_bazel_c_build,session_manager,sessiond)
	$(call copy_bazel_c_build,sctpd/src,sctpd)
	$(call copy_bazel_c_build,connection_tracker/src,connectiond)
	$(call copy_bazel_c_build,li_agent/src,liagentd)

build_oai: ## Build OAI
	$(call run_cmake, $(C_BUILD)/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(COMMON_FLAGS) $(OAI_NOTEST_FLAGS))

format_all:
	find $(MAGMA_ROOT)/orc8r/gateway/c/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h" -o -iname "*.hpp" \) -exec \
	/usr/bin/clang-format-11 --style='{BasedOnStyle: Google, IncludeBlocks: Preserve, SortIncludes: false}' -i {} \;
	find $(MAGMA_ROOT)/lte/gateway/c/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h"  -o -iname "*.hpp" \) -not \( -path "$(MAGMA_ROOT)/lte/gateway/c/core/*" -prune \) -exec \
	/usr/bin/clang-format-11 --style='{BasedOnStyle: Google, IncludeBlocks: Preserve, SortIncludes: false}' -i {} \;
	find $(MAGMA_ROOT)/lte/gateway/c/core/ \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h" -o -iname "*.hpp" \) -exec \
	clang-format --style=file -i {} \;

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

test_c: ## Run all Bazel-ified C/C++ tests
	bazel test $(BAZEL_FLAGS) -- //orc8r/gateway/c/...:* //lte/gateway/c/...:* -//lte/gateway/c/core/...:*

test_oai: ## Run all OAI-specific tests
	$(call run_ctest, $(C_BUILD)/core, $(C_BUILD)/core/oai, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(OAI_TEST_FLAGS), $(OAI_TESTS))

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

build_oai_clang: ## Build OAI with Clang, store compiler outputs to log
	$(call run_cmake, $(C_BUILD)/core, $(GATEWAY_C_DIR)/core, $(OAI_FLAGS) $(COMMON_FLAGS) $(OAI_NOTEST_FLAGS), CC="clang" CXX="clang++") 2>&1 | tee /tmp/clang-build.oai.log

# Upload Clang-Warning counts by type to Google Sheet via Google Survey
#  Graph available at https://docs.google.com/spreadsheets/d/1ndiIKJNI2IJZBwavnu1x_KwwvlQppnoYEOtYdihqiIQ/edit#gid=734064581
# TODO: Move to a third party service (I could not find) or make a binary that parses + uploads this without missing new warning types.
clang_warning_oai_upload: build_oai_clang ## Generate and then upload warning statistics to Google Survey for telemetry
	curl -L -v -G -Ss \
		--data-urlencode "entry.608561103=$(BRANCH)" \
		--data-urlencode "entry.1257088109=$(REVISION)" \
		--data-urlencode "entry.708918097=OAI" \
		--data-urlencode "entry.1995144005=$(shell grep '\[-Wconstant-conversion\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.889426184=$(shell grep '\[-Wdeprecated-declarations\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.6143542=$(shell grep '\[-Wenum-conversion\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1953446373=$(shell grep '\[-Wextern-c-compat\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1508747323=$(shell grep '\[-Winconsistent-missing-override\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1682903653=$(shell grep '\[-Winitializer-overrides\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1030484709=$(shell grep '\[-Wnon-c-typedef-for-linkage\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1885910120=$(shell grep '\[-Wnull-dereference\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.411691731=$(shell grep '\[-Wparentheses-equality\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.872392934=$(shell grep '\[-Wpointer-bool-conversion\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1058267165=$(shell grep '\[-Wtautological-constant-out-of-range-compare\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1145616398=$(shell grep '\[-Wtautological-overlap-compare\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.1265654844=$(shell grep '\[-Wtautological-pointer-compare\]' /tmp/clang-build.oai.log | wc -l)" \
		--data-urlencode "entry.599425517=$(shell grep '\[-Wtypedef-redefinition\]' /tmp/clang-build.oai.log | wc -l)" \
		https://docs.google.com/forms/d/e/1FAIpQLScKB3nLPASxzr4AXW5_yeHCjkEURY0K9OAFPIyNFzkA5CY_kw/formResponse?usp=pp_url

# Run clang-tidy
# TODO: CMake config issue - compile_commands.json is only generated if you first make_oai, then test_oai (or do either twice).
#       Additionally compile_commands.json is **not capturing all build artifacts this way**
clang_tidy_oai:
	mkdir -p $(C_BUILD)/core/oai/build
	cd $(C_BUILD)/core/oai/build;cmake $(GATEWAY_C_DIR)/core
	sed -i 's/CMAKE_EXPORT_COMPILE_COMMANDS:BOOL=OFF/CMAKE_EXPORT_COMPILE_COMMANDS:BOOL=ON/g' $(C_BUILD)/core/oai/build/CMakeCache.txt
	cmake --build $(C_BUILD)/core/oai/build/
	cd $(C_BUILD)/core/oai/build;wget https://raw.githubusercontent.com/llvm-mirror/clang-tools-extra/master/clang-tidy/tool/run-clang-tidy.py;\
	python run-clang-tidy.py -p $(C_BUILD)/core/oai/build/oai -j 2 -checks='-*,clang-analyzer-security*,android-*,cert-*,clang-analyzer-*,concurrency,misc-*,-misc-unused-parameters,bugprone-*' 2>&1 | tee /tmp/clang-tidy-oai.findings

# Pushes per-finding counts to:
#   https://docs.google.com/forms/d/1-45BZTHh4uBBOCqYM4LAD4zFT91B47sEEntHLkQuXWA/edit#responses
clang_tidy_oai_upload: clang_tidy_oai
	# Generate a summary of per-finding counts in the log
	grep '\]' /tmp/clang-tidy-oai.findings | grep warning: | awk -F'[][]' '{print $$2}' | sort | uniq -c
	curl -L -v -G -Ss \
		--data-urlencode "entry.338748281=$(BRANCH)" \
		--data-urlencode "entry.1421502557=$(REVISION)" \
		--data-urlencode "entry.1687500610=OAI" \
		--data-urlencode "entry.1690779794=$(shell grep '\[android-cloexec' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1926965966=$(shell grep '\[bugprone-branch-clone\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.951746576=$(shell grep '\[bugprone-macro-parentheses\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.593794544=$(shell grep '\[bugprone-narrowing-conversions\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.607048444=$(shell grep '\[bugprone-reserved-identifier,cert-dcl37-c,cert-dcl51-cpp\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.2081825272=$(shell grep '\[bugprone-signed-char-misuse,cert-str34-c\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.268408734=$(shell grep '\[bugprone-sizeof-expression\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1339652996=$(shell grep '\[bugprone-suspicious-string-compare\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1446914229=$(shell grep '\[bugprone-too-small-loop-variable\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.718845718=$(shell grep '\[bugprone-undefined-memory-manipulation\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.851497596=$(shell grep '\[cert-' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1963803956=$(shell grep '\[clang-analyzer-core.CallAndMessage\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1218571043=$(shell grep '\[clang-analyzer-core.NullDereference\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1618711162=$(shell grep '\[clang-analyzer-core.uninitialized.Assign\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1062688275=$(shell grep '\[clang-analyzer-deadcode.DeadStores\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1790897376=$(shell grep '\[clang-analyzer-optin' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.1730273925=$(shell grep '\[clang-analyzer-security.insecureAPI.strcpy\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.806421502=$(shell grep '\[clang-analyzer-unix.Malloc\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.999694738=$(shell grep '\[clang-analyzer-unix.MallocSizeof\]' /tmp/clang-tidy-oai.findings | wc -l)" \
		--data-urlencode "entry.740240539=$(shell grep '\[misc-' /tmp/clang-tidy-oai.findings | wc -l)" \
		https://docs.google.com/forms/d/e/1FAIpQLSfHGqOmDhUMAWHSjA_w6NOGqglQBx2IaO1bXLo6zrOE95sRWQ/formResponse?usp=pp_url

## Generate complete code structural information prior to any test execution
base_coverage: build_oai
	lcov --initial --directory $(C_BUILD) -c --output-file /tmp/coverage_initialize.info.raw
	lcov -r /tmp/coverage_initialize.info.raw "/*/test/*" "/usr/*" "/build/*protos*" -o /tmp/coverage_initialize.info
	rm -f `find $(C_BUILD) -name *.gcda` # Clean up any prior coverage data

# Combine results of sub-coverages
coverage: base_coverage coverage_oai
	lcov -a /tmp/coverage_initialize.info -a /tmp/coverage_oai.info -o $(C_BUILD)/coverage.info
