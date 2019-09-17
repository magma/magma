# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

list(APPEND CMAKE_CXX_FLAGS
  "-fdiagnostics-show-option"
  "-fstack-protector-all"
  "-g"
  "-ggdb3"
  "-O1"
  "-fno-omit-frame-pointer")
