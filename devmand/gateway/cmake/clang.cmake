# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}" "-std=c++2a")
list(APPEND CMAKE_CXX_WARNINGS
  "-fcolor-diagnostics"
  "-Weverything"
  "-Wno-weak-vtables"
  "-Wno-\#pragma-messages"
  "-Wno-unknown-attributes"
  "-Wno-address-of-packed-member"
  "-Wno-c++98-compat"
  "-Wno-c++98-compat-pedantic"
  "-Wno-padded"
  "-Wno-packed"
  "-Wno-disabled-macro-expansion"
  "-Weffc++"
  "-Werror"
  "-pedantic")
