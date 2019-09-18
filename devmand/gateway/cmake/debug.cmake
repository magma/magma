# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

list(APPEND CMAKE_CXX_FLAGS
  "-fdiagnostics-show-option"
  "-fstack-protector-all"
  "-g"
  "-ggdb3"
  "-O1"
  "-fno-omit-frame-pointer")
