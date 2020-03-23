# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

cmake_minimum_required(VERSION 3.0.2)

PROJECT(MagmaCommon)

include(CMakeProtoMacros.txt)

ADD_SUBDIRECTORY(service303)
ADD_SUBDIRECTORY(service_registry)
ADD_SUBDIRECTORY(async_grpc)
ADD_SUBDIRECTORY(config)
ADD_SUBDIRECTORY(datastore)
ADD_SUBDIRECTORY(policydb)
ADD_SUBDIRECTORY(logging)
ADD_SUBDIRECTORY(eventd)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)
