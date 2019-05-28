# Copyright (c) 2017-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")
add_definitions(-DDEBUG_IS_ON=1)

add_library(LIB_OPENFLOW_CONTROLLER
  ControllerMain.cpp
  OpenflowController.cpp
  PagingApplication.cpp
  ControllerEvents.cpp
  BaseApplication.cpp
  OpenflowMessenger.cpp
  GTPApplication.cpp
  IMSIEncoder.cpp
  )
target_link_libraries(LIB_OPENFLOW_CONTROLLER
  COMMON
  fluid_base fluid_msg
  LIB_BSTR
  TASK_SGW
)
target_include_directories(LIB_OPENFLOW_CONTROLLER PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)
