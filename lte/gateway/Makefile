.PHONY: all build test clean run status

GATEWAY_C_DIR = $(MAGMA_ROOT)/lte/gateway/c
GRPC_CPP_PLUGIN_PATH ?= `which grpc_cpp_plugin`
BUILD_TYPE ?= Debug


FEATURES ?= mme
# EXCLUSIVE_FEATURE_LIST : list of primary features that cannot be requested
# together.
EXCLUSIVE_FEATURE_LIST = mme agw_gtpnl agw_of
# AVAILABLE_FEATURE_LIST : every feature not in this list will trigger an error.
AVAILABLE_FEATURE_LIST =  mme  
REQUESTED_FEATURE_LIST = $(sort $(FEATURES))

OAI_FLAGS = -DEMBEDDED_SGW=False 
FUZZ_FLAGS = $(OAI_FLAGS) -DFUZZ=True
TEST_FLAG = -DBUILD_TESTS=1

all: build

build: build_python build_common build_oai build_sctpd build_session_manager

test: test_python test_oai test_session_manager

# Requires install of afl from http://lcamtuf.coredump.cx/afl/
fuzz: build_fuzz
	mkdir -p ~/fuzz/input ~/fuzz/output && \
	echo "testtest" > ~/fuzz/input/sample1 && \
	afl-fuzz -i ~/fuzz/input/ -o ~/fuzz/output/ -m 300M ~/build/c/oai_fuzz/oai_fuzz/oai_fuzz nas

clean: clean_python
	rm -rf $(C_BUILD)

clean_python:
	make -C $(MAGMA_ROOT)/lte/gateway/python clean

start:
	sudo service magma@magmad start

stop:
	sudo service magma@* stop

restart: stop start

run: build restart

status:
	sudo service magma@* status

# run_cmake BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS, ENV
define run_cmake
mkdir -p $(1)
cd $(1) && $(4) cmake $(2) -DCMAKE_BUILD_TYPE=$(BUILD_TYPE) $(3) -GNinja
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

# run_ctest BUILD_DIRECTORY, FILE_DIRECTORY, FLAGS
define run_ctest
$(call run_cmake, $(1), $(2), $(3) $(TEST_FLAG))
cd $(1) && ctest --output-on-failure
endef

build_python: stop
	make -C $(MAGMA_ROOT)/lte/gateway/python buildenv

build_common:
	$(call run_cmake, $(C_BUILD)/magma_common, $(MAGMA_ROOT)/orc8r/gateway/c/common, )

build_oai: build_common
	$(call run_cmake, $(C_BUILD)/oai, $(GATEWAY_C_DIR)/oai, $(OAI_FLAGS))

build_sctpd: build_common
	$(call run_cmake, $(C_BUILD)/sctpd, $(GATEWAY_C_DIR)/sctpd, )

# Requires install of afl from http://lcamtuf.coredump.cx/afl/
build_fuzz:
	$(call run_cmake, $(C_BUILD)/oai_fuzz, $(GATEWAY_C_DIR)/oai, $(FUZZ_FLAGS), CC=afl-gcc CXX=afl-g++ AFL_USE_ASAN=1)

build_session_manager: build_common
	$(call run_cmake, $(C_BUILD)/session_manager, $(GATEWAY_C_DIR)/session_manager, )

# Catch all for c services that don't have custom flags
# This works with build_dpi
build_%:
	$(call run_cmake, $(C_BUILD)/$*, $(MAGMA_ROOT)/c/$*, )

scan_oai:
	$(call run_scanbuild, $(C_BUILD)/scan/oai, $(GATEWAY_C_DIR)/oai, $(OAI_FLAGS))

format_oai:
	find $(GATEWAY_C_DIR)/oai \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h" \) -exec \
	clang-format --style=file -i {} \;

test_python: stop
	make -C $(MAGMA_ROOT)/lte/gateway/python test_all

test_oai: build_common
	$(call run_ctest, $(C_BUILD)/oai, $(GATEWAY_C_DIR)/oai, $(OAI_FLAGS))

# Catch all for c service tests
# This works with test_dpi and test_session_manager
test_%: stop build_common
	$(call run_ctest, $(C_BUILD)/$*, $(GATEWAY_C_DIR)/$*, )

# TODO: include coverage of dpim when it is included in make build
COV_OUTPUT_OAI = $(OAI_BUILD)/coverage.info
COV_OUTPUT_SM = $(C_BUILD)/session_manager/coverage.info
COV_OUTPUT_TOTAL = $(C_BUILD)/coverage.info
# Put HTML within magma directory, so it can be accessed from outside VM
COV_HTML_DIR_OAI = $(GATEWAY_C_DIR)/oai/code_coverage
COV_HTML_DIR_SM = $(GATEWAY_C_DIR)/session_manager/code_coverage
COV_HTML_DIR_TOTAL = $(MAGMA_ROOT)/c/code_coverage

SLEEP_SECS = 10

coverage:
	# Remove any previous gcov output files
	rm -f `find $(OAI_BUILD) -name *.gcda`
	rm -f `find $(C_BUILD)/session_manager -name *.gcda`
	# Stop processes to generate gcda
	sudo pkill -USR1 mme
	sudo pkill -INT sessiond
	# Wait for gcov output to be generated
	@echo "Waiting for $(SLEEP_SECS) for gcov to write files"
	sleep $(SLEEP_SECS)
	# Capture coverage info under different directories individually
	lcov --capture --directory $(OAI_BUILD) --output-file $(COV_OUTPUT_OAI)
	lcov --capture --directory $(C_BUILD)/session_manager --output-file $(COV_OUTPUT_SM)
	# Merge coverage info
	lcov --add-tracefile $(COV_OUTPUT_OAI) --add-tracefile $(COV_OUTPUT_SM) --output-file $(COV_OUTPUT_TOTAL)
	# Remove coverage info regarding libraries from /usr/include/
	lcov --remove $(COV_OUTPUT_TOTAL) '/usr/include/*' '/usr/local/include/*' -o $(COV_OUTPUT_TOTAL) --quiet
	# Generate html which shows coverage with graph
	genhtml $(COV_OUTPUT_TOTAL) --output-directory $(COV_HTML_DIR_TOTAL)
	@echo "Generated coverage output to $(COV_HTML_DIR_TOTAL)/index.html"

coverage_oai:
	rm -f `find $(OAI_BUILD) -name *.gcda`
	sudo pkill -USR1 mme
	@echo "Waiting for $(SLEEP_SECS) for gcov to write files"
	sleep $(SLEEP_SECS)
	lcov --capture --directory $(OAI_BUILD) --output-file $(COV_OUTPUT_OAI)
	lcov --remove $(COV_OUTPUT_OAI) '/usr/include/*' '/usr/local/include/*' -o $(COV_OUTPUT_OAI) --quiet
	genhtml $(COV_OUTPUT_OAI) --output-directory $(COV_HTML_DIR_OAI)
	@echo "Generated coverage output to $(COV_HTML_DIR_OAI)/index.html"

coverage_sm:
	rm -f `find $(C_BUILD)/session_manager -name *.gcda`
	sudo pkill -INT sessiond
	sleep $(SLEEP_SECS)
	lcov --capture --directory $(C_BUILD)/session_manager --output-file $(COV_OUTPUT_SM)
	lcov --remove $(COV_OUTPUT_SM) '/usr/include/*' '/usr/local/include/*' -o $(COV_OUTPUT_SM) --quiet
	genhtml $(COV_OUTPUT_SM) --output-directory $(COV_HTML_DIR_SM)
	@echo "Generated coverage output to $(COV_HTML_DIR_SM)/index.html"

code_stats:
	# Get line counts for entire magma project. Exclude auto-generated code.
	sudo apt-get install -y cloc
	cloc . --exclude-dir=cloud/go/src/magma/obsidian/models
