#!/bin/bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#

# shellcheck disable=SC1091
source /etc/environment

cd "$CODEGEN_ROOT" && "$M2_HOME"/bin/mvn clean package && cd - || exit
