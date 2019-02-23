# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
ifndef PROTO_LIST
PROTO_LIST:=orc8r_protos
endif

TESTS=magma/common/tests \
      magma/configuration/tests \
      magma/magmad/kernel_check/tests \
      magma/magmad/network_check/tests \
      magma/magmad/tests \
      magma/magmad/logging/tests \
      magma/magmad/upgrade/tests \
      magma/metricsd/tests
