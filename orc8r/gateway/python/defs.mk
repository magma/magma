# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
ifndef PROTO_LIST
PROTO_LIST:=orc8r_protos
endif

TESTS=magma/common/redis/tests \
      magma/common/tests \
      magma/configuration/tests \
      magma/directoryd/tests \
      magma/eventd/tests \
      magma/magmad/check/tests \
      magma/magmad/check/kernel_check/tests \
      magma/magmad/check/machine_check/tests \
      magma/magmad/check/network_check/tests \
      magma/magmad/tests \
      magma/magmad/logging/tests \
      magma/magmad/upgrade/tests \
      magma/magmad/generic_command/tests \
      magma/state/tests \
      magma/tests
