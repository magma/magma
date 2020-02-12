# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
PYTHON_SRCS=$(MAGMA_ROOT)/lte/gateway/python $(MAGMA_ROOT)/orc8r/gateway/python
PROTO_LIST:=orc8r_protos lte_protos feg_protos
SWAGGER_LIST:=orc8r_swagger_specs

# Path to the test files
TESTS=magma/tests \
	  magma/policydb/tests \
	  magma/enodebd/tests \
      magma/mobilityd/tests \
      magma/pipelined/openflow/tests \
      magma/pkt_tester/tests \
      magma/redirectd/tests \
      magma/subscriberdb/tests

SUDO_TESTS=magma/pipelined/tests
