#!/bin/bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#

# open config has some non-standard complaint yang models. This script takes in
# a model and makes it compliant.

function compliance {
  local pattern="$1"
  sed -i'' "${pattern}" "${model}"
}

model="$1"

# no when statement on state data
compliance 's;when "oc-if:state/oc-if:type;when "oc-if:config/oc-if:type;g'

# Use wc3 regex not posix
compliance "s;\(\s'\)\^;\1;g"
compliance "s;\$\('\;\);\1;g"
