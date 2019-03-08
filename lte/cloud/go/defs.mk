# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
$(MODULE)_ROOTLEVEL_PKG := magma/lte/cloud/go

$(MODULE)_SERVICE_NAMES := directoryd subscriberdb eps_authentication policydb meteringd_records

$(MODULE)_SWAGGER_DIRS := services/cellular/swagger services/policydb/swagger services/subscriberdb/swagger services/meteringd_records/swagger
$(MODULE)_SWAGGER_GEN := definitions_cellular.yml:services/cellular/obsidian definitions_policies.yml:services/policydb/obsidian definitions_subscribers.yml:services/subscriberdb/obsidian definitions_metering.yml:services/meteringd_records/obsidian
