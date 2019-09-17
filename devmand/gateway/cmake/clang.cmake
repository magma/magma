# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
