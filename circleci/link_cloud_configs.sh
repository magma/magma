#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#

module_array=(`echo ${MAGMA_MODULES}`);
for mod in "${module_array[@]}"; do
    sudo ln -s ${mod}/cloud/configs /etc/magma/configs/$(basename ${mod})
done
