# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

project (devmand CXX)

cmake_minimum_required(VERSION 3.7.2)

set(MGM_COM_DIR /usr/lib/)
add_definitions(-DPROTOBUF_INLINE_NOT_IN_HEADERS=0)

include_directories("/usr/include/prometheus")
include_directories(${PROJECT_BINARY_DIR})
include_directories(${PROJECT_SOURCE_DIR}/src)

find_library(NETSNMPAGENT netsnmpagent)
find_library(NETSNMPMIBS netsnmpmibs)
find_library(NETSNMP netsnmp)
find_library(MGM_ASYNC_GRPC ASYNC_GRPC ${MGM_COM_DIR}/async_grpc)
find_library(MGM_CONFIG CONFIG HINTS ${MGM_COM_DIR}/config)
find_library(MGM_SERVICE303_LIB SERVICE303_LIB HINTS ${MGM_COM_DIR}/service303)
find_library(MGM_DATASTORE DATASTORE HINTS ${MGM_COM_DIR}/datastore)
find_library(MGM_POLICYDB POLICYDB HINTS ${MGM_COM_DIR}/policydb)
find_library(MGM_SRVC_REG SERVICE_REGISTRY HINTS ${MGM_COM_DIR}/service_registry)

add_library(devman
  ${PROJECT_SOURCE_DIR}/src/devmand/Application.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/cartography/Cartographer.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/cartography/Method.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/http/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/http/Response.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cnmaestro/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/mikrotik/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/mikrotik/WriteTask.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/packet/Engine.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/Config.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/CambiumDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/DemoDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/DcsgDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Factory.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/EchoDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/FrinxDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/mikrotik/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/mikrotik/Mib.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/State.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/DhcpdConfig.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/ErrorHandler.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/ErrorQueue.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/FileUtils.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/FileWatcher.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/device/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/interface/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/wifi/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/YangUtils.cpp
)

target_link_libraries(devman
  boost_fiber
  boost_context
  boost_thread
  boost_system
  crypto
  folly
  pthread
  double-conversion
  stdc++
  stdc++fs
  dl
  m
  event
  event_core
  ssh
  ssh_threads
  pcre
  ydk
  curl
  xslt
  xml2
  yaml-cpp
  protobuf
  cpp_redis
  tacopie
  glog
  gflags
  iberty
  unwind)

add_library(devman_channel_snmp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Engine.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/IfMib.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Oid.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Pdu.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Request.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Response.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Snmpv2Device.cpp)

target_link_libraries(devman_channel_snmp
  devman
  ${NETSNMPAGENT}
  ${NETSNMPMIBS}
  ${NETSNMP})

add_library(devman_service_magma
  ${PROJECT_SOURCE_DIR}/src/devmand/magma/DevConf.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/magma/Service.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/Service.cpp)

target_link_libraries(devman_service_magma
  boost_fiber
  boost_context
  folly
  pthread
  double-conversion
  stdc++
  dl
  m
  event
  event_core
  ${MGM_SERVICE303_LIB}
  ${MGM_SRVC_REG}
  ${MGM_ASYNC_GRPC}
  ${MGM_CONFIG}
  ${MGM_POLICYDB}
  ${MGM_DATASTORE}
  grpc++
  grpc
  prometheus-cpp
  yaml-cpp
  protobuf
  cpp_redis
  tacopie
  glog
  gflags
  iberty
  unwind)

add_executable(devmand
  ${PROJECT_SOURCE_DIR}/src/devmand/main.cpp)

target_link_libraries(devmand
  devman
  devman_channel_snmp
  devman_service_magma)

add_executable(devmantest
  ${PROJECT_SOURCE_DIR}/src/devmand/test/MikrotikChannelTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/DemoDeviceTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/DevConfTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/EventBaseTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/FileWatcherTest.cpp)

target_link_libraries(devmantest
  devman
  devman_channel_snmp
  gtest
  gtest_main
  devman_service_magma)

install(TARGETS devmand devmantest
  RUNTIME DESTINATION bin)

include (CTest)

add_test(devmantest devmantest)

set(CMAKE_CXX_COMPILER gcc)

if(CMAKE_CXX_COMPILER MATCHES ".*clang.*")
  include(cmake/clang.cmake)
elseif(CMAKE_CXX_COMPILER MATCHES ".*(gcc|g\\+\\+).*")
  include(cmake/gcc.cmake)
endif()

include(cmake/${CMAKE_BUILD_TYPE}.cmake)

list(APPEND CMAKE_CXX_FLAGS "${CMAKE_CXX_WARNINGS}")
string(REPLACE ";" " " CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}")
