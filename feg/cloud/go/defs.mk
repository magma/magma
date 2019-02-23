# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
$(MODULE)_ROOTLEVEL_PKG := magma/feg/cloud/go

$(MODULE)_SERVICE_NAMES := feg_relay health

$(MODULE)_SWAGGER_DIRS := services/controller/swagger
$(MODULE)_SWAGGER_GEN := definitions_feg.yml:services/controller/obsidian
