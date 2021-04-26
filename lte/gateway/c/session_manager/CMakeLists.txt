cmake_minimum_required(VERSION 3.0.2)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_VERBOSE_MAKEFILE ON)

project(MagmaSessionManager)
include(ExternalProject)

ExternalProject_Add(MagmaCommon
    SOURCE_DIR "$ENV{MAGMA_ROOT}/orc8r/gateway/c/common"
    BINARY_DIR "${CMAKE_BINARY_DIR}/common"
    INSTALL_COMMAND "")

ExternalProject_Add(SessionManager
    SOURCE_DIR "${CMAKE_SOURCE_DIR}/src"
    BINARY_DIR "${CMAKE_BINARY_DIR}/src"
    INSTALL_COMMAND ""
    DEPENDS MagmaCommon)