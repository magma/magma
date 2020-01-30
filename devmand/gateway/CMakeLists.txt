# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

project (devmand CXX)

cmake_minimum_required(VERSION 3.7.2)

set(MGM_COM_DIR /usr/lib/)
# Pass libyang plugins dir for YDK libyang
add_definitions(-DLIBYANG_PLUGINS_DIR="/usr/lib/libyang/")

# Pass libyang extensions and user_types for non-YDK libyang
add_definitions(-DLLLIBYANG_EXTENSIONS_PLUGINS_DIR="/usr/lib/libyang/extensions")
add_definitions(-DLIBYANG_USER_TYPES_PLUGINS_DIR="/usr/lib/libyang/user_types")

add_definitions(-DPROTOBUF_INLINE_NOT_IN_HEADERS=0)

include_directories("/usr/include/prometheus")
include_directories(${PROJECT_BINARY_DIR})
include_directories(${PROJECT_SOURCE_DIR}/src)

find_library(NETSNMPAGENT netsnmpagent)
find_library(NETSNMPMIBS netsnmpmibs)
find_library(NETSNMP netsnmp)

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
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/ping/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/ping/Engine.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/Config.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cambium/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Datastore.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/demo/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/Factory.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/echo/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/frinx/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/ping/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/mikrotik/Device.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/mikrotik/Mib.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/error/ErrorHandler.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/error/ErrorQueue.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/MetricSink.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/device/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/interface/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/models/wifi/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/syslog/Manager.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/utils/ConfigGenerator.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/utils/FileUtils.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/utils/FileWatcher.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/utils/Time.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/utils/YangUtils.cpp
)

target_link_libraries(devman
  boost_context
  boost_thread
  boost_system
  boost_filesystem
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
  event_pthreads
  ssh
  ssh_threads
  pcre
  yang
  ydk
  ydk_openconfig
  ydk_ietf
  curl
  xslt
  xml2
  yaml-cpp
  protobuf
  cpp_redis
  tacopie
  glog
  gflags
  iberty)

add_library(devman_channel_snmp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Engine.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/EventHandler.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/IfMib.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Oid.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Pdu.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Request.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/snmp/Response.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/snmpv2/Device.cpp)

target_link_libraries(devman_channel_snmp
  devman
  ${NETSNMPAGENT}
  ${NETSNMPMIBS}
  ${NETSNMP})

add_library(devman_channel_cli
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/Command.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/engine/Engine.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/SshSession.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/Channel.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/PromptAwareCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/QueuedCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/SshSessionAsync.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/CliFlavour.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/SshSessionAsync.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/PlaintextCliDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/ParsingUtils.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/StructuredUbntDevice.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/schema/ModelRegistry.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/schema/Model.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/schema/Path.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/schema/SchemaContext.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/devices/cli/schema/BindingContext.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/Spd2Glog.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/IoConfigurationBuilder.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/ReadCachingCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/SshSocketReader.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/TimeoutTrackingCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/KeepaliveCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/ReconnectingCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/CliThreadWheelTimekeeper.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/CliTimekeeperWrapper.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/LoggingCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/CancelableWTCallback.h
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/TreeCacheCli.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/channels/cli/TreeCache.cpp
  )

target_link_libraries(devman_channel_cli
  devman)

find_library(MGM_ASYNC_GRPC ASYNC_GRPC ${MGM_COM_DIR}/async_grpc)
find_library(MGM_CONFIG CONFIG HINTS ${MGM_COM_DIR}/config)
find_library(MGM_SERVICE303_LIB SERVICE303_LIB HINTS ${MGM_COM_DIR}/service303)
find_library(MGM_DATASTORE DATASTORE HINTS ${MGM_COM_DIR}/datastore)
find_library(MGM_POLICYDB POLICYDB HINTS ${MGM_COM_DIR}/policydb)
find_library(MGM_SRVC_REG SERVICE_REGISTRY HINTS ${MGM_COM_DIR}/service_registry)

add_library(devman_service
  ${PROJECT_SOURCE_DIR}/src/devmand/Service.cpp)

target_link_libraries(devman_service
  folly
  pthread
  double-conversion
  stdc++)

add_library(devman_service_fscache
  ${PROJECT_SOURCE_DIR}/src/devmand/fscache/Service.cpp)

target_link_libraries(devman_service_fscache
  folly
  pthread
  double-conversion
  devman_service
  stdc++
  dl
  m
  event
  event_core
  prometheus-cpp
  yaml-cpp
  protobuf
  cpp_redis
  tacopie
  glog
  gflags
  iberty)

add_library(devman_service_magma
  ${PROJECT_SOURCE_DIR}/src/devmand/magma/DevConf.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/magma/Service.cpp)

target_link_libraries(devman_service_magma
  folly
  pthread
  devman_service
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
  iberty)

add_executable(devmand
  ${PROJECT_SOURCE_DIR}/src/devmand/main.cpp)

target_link_libraries(devmand
  devman
  devman_channel_snmp
  devman_channel_cli
  devman_service_fscache
  devman_service_magma)

add_executable(devmantest
  ${PROJECT_SOURCE_DIR}/src/devmand/test/ConfigGeneratorTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/RealCliDeviceTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/CliScaleTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/CliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/PromptAwareCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/TimeoutTrackingCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/SshSessionTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/QueuedCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/TimeoutTrackingCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/PromptAwareCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/ReadCachingCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/KeepaliveCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/ReconnectingCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/utils/Json.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/utils/Ssh.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/utils/Log.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/CommandTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/PathTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/SchemaContextTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/BindingContextTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/PlaintextCliDeviceTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/StructuredUbntDeviceTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/DemoDeviceTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/DevConfTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/SyslogManagerTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/DiffTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/ErrorQueueTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/EventBaseTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/FileWatcherTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/MikrotikChannelTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/PingChannelTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/SnmpChannelTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/ReconnectingSshTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/TreeCacheCliTest.cpp
  ${PROJECT_SOURCE_DIR}/src/devmand/test/cli/TreeCacheTest.cpp
)

target_link_libraries(devmantest
  devman
  devman_channel_snmp
  devman_channel_cli
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
