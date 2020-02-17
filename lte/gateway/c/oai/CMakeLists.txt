cmake_minimum_required(VERSION 3.0.2)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

project(MagmaCore)

include(CMakeOptionsMacros.txt)
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

################################################################
# Build type
################################################################
if (CMAKE_BUILD_TYPE STREQUAL "")
  set(CMAKE_BUILD_TYPE "RelWithDebInfo")
endif ()

add_list_string_option(CMAKE_BUILD_TYPE "RelWithDebInfo" "Choose the type of build, options are: None(CMAKE_CXX_FLAGS or CMAKE_C_FLAGS used) Debug Release RelWithDebInfo MinSizeRel." Debug Release RelWithDebInfo MinSizeRel)
message("Build type is ${CMAKE_BUILD_TYPE}")

if (CMAKE_BUILD_TYPE STREQUAL "Debug")
  add_boolean_option(LINK_GCOV True "Whether to link gcov")
endif ()

################################################################
# Compiler options
################################################################

# TODO (amar) check all the flags used. Copied as is from OAI.
set(CMAKE_C_FLAGS
  "${CMAKE_C_FLAGS} ${C_FLAGS_PROCESSOR} -std=gnu99 -Wall -Werror -Wstrict-prototypes -fno-strict-aliasing -rdynamic -funroll-loops -Wno-packed-bitfield-compat -fPIC"
)

# add autoTOOLS definitions that were maybe used
set(CMAKE_C_FLAGS
  "${CMAKE_C_FLAGS} -DSTDC_HEADERS=1 -DHAVE_SYS_TYPES_H=1 -DHAVE_SYS_STAT_H=1 -DHAVE_STDLIB_H=1 -DHAVE_STRING_H=1 -DHAVE_MEMORY_H=1 -DHAVE_STRINGS_H=1 -DHAVE_INTTYPES_H=1 -DHAVE_STDINT_H=1 -DHAVE_UNISTD_H=1 -DHAVE_FCNTL_H=1 -DHAVE_ARPA_INET_H=1 -DHAVE_SYS_TIME_H=1 -DHAVE_SYS_SOCKET_H=1 -DHAVE_STRERROR=1 -DHAVE_SOCKET=1 -DHAVE_MEMSET=1 -DHAVE_GETTIMEOFDAY=1 -DHAVE_STDLIB_H=1 -DHAVE_MALLOC=1 -DHAVE_LIBSCTP"
)

# set a flag for changes in the source code
# these changes are related to hardcoded path to include .h files
add_definitions(-DCMAKER)

# Generate *.gcno rather than *.c.gcno files for gcov
set(CMAKE_C_OUTPUT_EXTENSION_REPLACE 1)

if (NOT FUZZ)
  set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS} -fstack-protector-all -g -DMALLOC_CHECK_=3 -DDEBUG_IS_ON=1 -DTRACE_IS_ON=1 -O0 -fsanitize=address -fno-omit-frame-pointer")
  set(CMAKE_C_FLAGS_RELWITHDEBINFO "${CMAKE_C_FLAGS} -g -DMALLOC_CHECK_=3  -DDEBUG_IS_ON -O1")
  set(CMAKE_EXE_LINKER_FLAGS_DEBUG  "${CMAKE_EXE_LINKER_FLAGS_DEBUG} -fsanitize=address")

  set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -DDEBUG_IS_ON=1")
  set(CMAKE_CXX_FLAGS_RELWITHDEBINFO "${CMAKE_CXX_FLAGS_DEBUG} -DDEBUG_IS_ON=1")

  # Add Gcov flags
  set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS_DEBUG} -g -O0 -fprofile-arcs -ftest-coverage")
  # Add gcov linker flags
  set(CMAKE_EXE_LINKER_FLAGS_DEBUG  "${CMAKE_EXE_LINKER_FLAGS_DEBUG} -fprofile-arcs -ftest-coverage")

  # Add LeakSanitizer (lsan) support to the release build of MME so that we can
  # find memory leaks in production.
  set(CMAKE_C_FLAGS_RELWITHDEBINFO "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
  set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO  "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")
endif (NOT FUZZ)

################################################################
# Set proto directory for state proto definitions
################################################################
set(STATE_PROTO_DIR ${PROJECT_SOURCE_DIR}/protos)

################################################################
# S1AP LAYER OPTIONS
################################################################
add_boolean_option(S1AP_DEBUG_LIST                 False    "Traces, option to be removed soon")

################################################################
# SCTP LAYER OPTIONS
################################################################
add_boolean_option(SCTP_DUMP_LIST                  False    "Traces, option to be removed soon")
add_boolean_option(TRACE_HASHTABLE                 False    "Trace hashtables operations ")
add_boolean_option(TRACE_3GPP_SPEC                 True     "Log hits of 3GPP specifications requirements")
add_boolean_option(ENABLE_OPENFLOW                 False    "Openflow based dataplane")
add_boolean_option(EMBEDDED_SGW                    False    "Add the SPGW task to the MME binary")
add_boolean_option(LINK_GCOV                       False    "Whether to link gcov")

if (EMBEDDED_SGW)
include(CMakeAgwOptions.txt)
else (EMBEDDED_SGW)
include(CMakeMmeOptions.txt)
endif (EMBEDDED_SGW)

################################################################
# Include CMake modules to find other library
################################################################

set(CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/modules/")

include(FindPkgConfig)
include(FindThreads)

pkg_search_module(CONFIG libconfig REQUIRED)
include_directories(${CONFIG_INCLUDE_DIRS})

set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
find_library(ASYNC_GRPC ASYNC_GRPC ${MAGMA_LIB_DIR}/async_grpc)
find_library(SERVICE303_LIB SERVICE303_LIB ${MAGMA_LIB_DIR}/service303)
find_library(CONFIG CONFIG ${MAGMA_LIB_DIR}/config)
find_library(SERVICE_REGISTRY SERVICE_REGISTRY ${MAGMA_LIB_DIR}/service_registry)

################################################################
# Add sub modules
################################################################

include_directories(${PROJECT_SOURCE_DIR}/include)

#############################################################
# Add magma libraries
#############################################################

add_subdirectory(common)
add_subdirectory(lib)
add_subdirectory(oai_mme)
add_subdirectory(oai_fuzz)
add_subdirectory(tasks)

if (BUILD_TESTS)
  enable_testing()
  add_subdirectory(test)
endif (BUILD_TESTS)
