#!/bin/sh
################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

TEMPLATE="${TEMPLATE_ENV:-radius.cwf.config.json.template}"

/usr/bin/envsubst < "${TEMPLATE}" > ./radius.config.json
./radius
