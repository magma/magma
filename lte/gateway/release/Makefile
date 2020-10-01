################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################
MAGMA_ROOT = ~/magma
PY_LTE     = $(MAGMA_ROOT)/lte/gateway/python
PY_ORC8R   = $(MAGMA_ROOT)/orc8r/gateway/python

magma.lockfile: $(PY_LTE)/setup.py $(PY_ORC8R)/setup.py
	./pydep finddep --install-from-official -l ./magma.lockfile $(PY_ORC8R)/setup.py $(PY_LTE)/setup.py
